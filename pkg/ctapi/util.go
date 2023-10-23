package ctapi

import (
	"unsafe"
)

func StringToLPCTSTR(value string) uintptr {
	bytes := append([]byte(value), 0)
	return uintptr(unsafe.Pointer(&bytes[0]))
}
