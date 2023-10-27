package ctapi

import (
	"golang.org/x/sys/windows"
	"strconv"
	"unsafe"
)

type CtList struct {
	hList windows.Handle
	api   *CtApi
}

func (this *CtApi) NewList(connection windows.Handle) (*CtList, error) {

	r1, _, err := this.procs.ctListNew.Call(uintptr(connection),
		uintptr(CT_LIST_EVENT))
	if r1 == 0 {
		return nil, err
	}

	return &CtList{
		api:   this,
		hList: windows.Handle(r1),
	}, nil
}

// Free releases a list back to CtApi.  Returns nil on success, otherwise
// returns an error.
func (this *CtList) Free() error {
	r1, _, err := this.api.procs.ctListFree.Call(uintptr(this.hList))
	if r1 == 0 {
		return err
	}
	return nil
}

// Add adds a tag to the list of tags being observed.  Returns
// a handle that can be used to read the item or remove it from
// the list.  Also returns an error code that is nil if there
// is no error
func (this *CtList) Add(sTag string) (windows.Handle, error) {
	r1, _, err := this.api.procs.ctListAdd.Call(uintptr(this.hList),
		StringToLPCTSTR(sTag))
	if r1 == 0 {
		return 0, err
	}
	return windows.Handle(r1), nil
}

func (this *CtList) Read() error {
	r1, _, err := this.api.procs.ctListRead.Call(uintptr(this.hList),
		nullptr)
	if r1 == 0 {
		return err
	}
	return nil
}

func (this *CtList) GetFloatValue(hTag windows.Handle) (float32, error) {
	const buflen = 255
	buf := make([]uint8, buflen)
	bufPtr := unsafe.Pointer(&buf[0])

	r1, _, err := this.api.procs.ctListData.Call(
		uintptr(hTag),   // hTag (HANDLE)
		uintptr(bufPtr), // pBuffer (void *)
		buflen,          // dwLength (DWORD)
		0)               // dwMode (set to zero)

	// returns false (zero) on error
	if r1 == 0 {
		return 0, err
	}

	strValue := windows.BytePtrToString(&buf[0])
	floatValue, err := strconv.ParseFloat(strValue, 32)

	return float32(floatValue), err
}
