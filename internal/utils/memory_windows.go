//go:build windows
// +build windows

package utils

import (
	"syscall"
	"unsafe"
)

func GetFreeSpace(path string) (uint64, error) {
    lpFreeBytesAvailable := int64(0)
    lpTotalNumberOfBytes := int64(0)
    lpTotalNumberOfFreeBytes := int64(0)

    pathPtr, _ := syscall.UTF16PtrFromString(path)
    r, _, err := syscall.NewLazyDLL("kernel32.dll").
        NewProc("GetDiskFreeSpaceExW").
        Call(uintptr(unsafe.Pointer(pathPtr)),
            uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
            uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
            uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)),
        )
    if r == 0 {
        return 0, err
    }
    return uint64(lpFreeBytesAvailable), nil
}
