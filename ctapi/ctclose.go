package ctapi

import "golang.org/x/sys/windows"

func (this *CtApi) CtClose(handle windows.Handle) (bool, error) {

	r1, _, err := this.procs.ctClose.Call(uintptr(handle))

	if r1 == 1 {
		return true, nil
	} else {
		return false, err
	}

}
