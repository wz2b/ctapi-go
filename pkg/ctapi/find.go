package ctapi

import (
	"fmt"
	"golang.org/x/sys/windows"
	"unsafe"
)

func (this *CtApi) CtFindFirst(h windows.Handle, table string) (windows.Handle, windows.Handle, error) {
	var hObj windows.Handle

	tblRef := StringToLPCTSTR(table)

	r1, _, err := this.procs.ctFindFirst.Call(
		uintptr(h),                     // hCTAPI (HANDLE)
		tblRef,                         // szTableName (LPCTSTR)
		nullptr,                        // szFilter (LPCTSTR)
		uintptr(unsafe.Pointer(&hObj)), // pObjHnd (HANDLE)
		0)                              // dwFlags (argument no longer used, pass 0)

	if r1 == 0 {
		return 0, 0, err
	}
	return windows.Handle(r1), hObj, nil
}

func (this *CtApi) CtFindNext(hSearch windows.Handle) (windows.Handle, error) {
	var hObj windows.Handle
	r1, _, err := this.procs.ctFindNext.Call(uintptr(hSearch), uintptr(unsafe.Pointer(&hObj)))

	if r1 == 1 {
		return hObj, nil
	} else {
		return 0, err
	}
}

func (this *CtApi) CtFindClose(h windows.Handle) error {
	r1, _, err := this.procs.ctFindClose.Call(uintptr(h), uintptr(h))
	if r1 == 1 {
		return nil
	} else {
		return err
	}
}

func (this *CtApi) FindAll(h windows.Handle, table string) <-chan windows.Handle {
	ch := make(chan windows.Handle)
	go func() {
		hSearch, firstObject, _ := this.CtFindFirst(h, table)
		defer this.CtFindClose(hSearch)

		hObj := firstObject
		if hSearch == 0 {
			println("No tags found")
			// no tags found
		} else {
			var err error
			for hObj != 0 {
				ch <- hObj
				hObj, err = this.CtFindNext(hSearch)
				if err != nil {
					fmt.Errorf("error geting next", err)
				}
			}
		}
		close(ch)
	}()

	return ch
}
