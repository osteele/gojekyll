//go:build windows

package utils

import "syscall"

// isWindowsDirNotEmpty checks if the error is Windows ERROR_DIR_NOT_EMPTY
func isWindowsDirNotEmpty(errno syscall.Errno) bool {
	return errno == syscall.ERROR_DIR_NOT_EMPTY
}
