//go:build darwin
// +build darwin

package utils

import "golang.org/x/sys/unix"

func GetFreeSpace(path string) (uint64, error) {
    var stat unix.Statfs_t
    if err := unix.Statfs(path, &stat); err != nil {
        return 0, err
    }
    return stat.Bavail * uint64(stat.Bsize), nil
}
