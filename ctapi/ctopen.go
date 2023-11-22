package ctapi

import (
	"golang.org/x/sys/windows"
)

func (this *CtApi) CtOpen() (windows.Handle, error) {
	r1, _, err := this.procs.ctOpen.Call(
		nullptr, nullptr, nullptr,
		CT_OPEN_READ_ONLY|CT_OPEN_BATCH|CT_OPEN_EXTENDED|CT_OPEN_WINDOWSUSER)

	if r1 != 0 {
		return windows.Handle(r1), nil
	}

	if err != nil {
		return 0, err
	}

	return windows.Handle(r1), nil
}
