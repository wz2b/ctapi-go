package ctapi

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"reflect"
	"unsafe"
)

func Init(dllName string, dllPath string) (*CtApi, error) {

	kernel, _ := windows.LoadDLL("kernel32.dll")
	defer kernel.Release()

	addDllDir, err := kernel.FindProc("AddDllDirectory")
	if err != nil {
		return nil, err
	}

	dllPathC, err := windows.UTF16PtrFromString(dllPath)
	if err != nil {
		panic(err)
	}

	r1, _, err := addDllDir.Call(uintptr(unsafe.Pointer(dllPathC)))
	if r1 == 0 {
		fmt.Println(os.Stderr, "Unable to add DLL path")
		return nil, err
	}

	h, err := windows.LoadLibraryEx(dllName, 0,
		windows.LOAD_LIBRARY_SEARCH_DEFAULT_DIRS|
			windows.LOAD_LIBRARY_SEARCH_USER_DIRS)

	dll := &windows.DLL{Name: dllName, Handle: h}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load DLL")
		return nil, err
	}

	ctSetManagedBinDirectory, err := dll.FindProc("ctSetManagedBinDirectory")
	if err != nil {
		panic(err)
	}
	fmt.Printf("ctSetManagedBinDirectory(): %v\n", ctSetManagedBinDirectory)

	fmt.Println("Setting managed directory path")
	tf, _, err := ctSetManagedBinDirectory.Call(StringToLPCTSTR(dllPath))
	if tf != 1 {
		fmt.Fprintln(os.Stderr, "Failed to set CTAPI managed bin directory")
		return nil, err
	}

	procs, err := getProcs(dll)
	if err != nil {
		return nil, err
	}

	return &CtApi{
		dll:   dll,
		procs: procs,
	}, nil
}

func getProcs(dll *windows.DLL) (*ctApiProcs, error) {
	var links ctApiProcs

	t := reflect.TypeOf(links)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		f, err := dll.FindProc(name)
		if err != nil {
			return nil, err
		}
		setFieldValue(&links, name, f)

	}

	fmt.Println("All functions linked")
	return &links, nil
}

func setFieldValue(target any, fieldName string, value any) {
	rv := reflect.ValueOf(target)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	if !rv.CanAddr() {
		panic("target must be addressable")
	}
	if rv.Kind() != reflect.Struct {
		panic(fmt.Sprintf(
			"unable to set the '%s' field value of the type %T, target must be a struct",
			fieldName,
			target,
		))
	}
	rf := rv.FieldByName(fieldName)

	reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}
