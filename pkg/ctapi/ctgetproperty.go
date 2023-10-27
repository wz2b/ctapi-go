package ctapi

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func (this *CtApi) GetPropertyAsString(hTag windows.Handle, name string) (string, error) {
	const buflen = 255
	//namePtr := unsafe.Pointer(StringToWideCharPtr(name))
	namePtr := unsafe.Pointer(StringToLPCTSTR(name))

	buf := make([]uint8, buflen)
	bufPtr := unsafe.Pointer(&buf[0])

	var resultLen uint32
	lenPtr := unsafe.Pointer(&resultLen)

	r1, _, err := this.procs.ctGetProperty.Call(
		uintptr(hTag),    // hnd (HANDLE)
		uintptr(namePtr), // szName (LPCTSTR)
		uintptr(bufPtr),  // pData (void *)
		buflen,           // eBufferLength (DWORD)
		uintptr(lenPtr),  // dwResultLength (DWORD *)
		DBTYPE_STR)       // dwType (DWORD)

	if r1 == 0 {
		return "ERR", err
	} else {
		return windows.BytePtrToString(&buf[0]), nil
	}
}
