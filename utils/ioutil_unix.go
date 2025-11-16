//go:build !windows

package utils

import "syscall"

// isWindowsDirNotEmpty always returns false on non-Windows platforms
func isWindowsDirNotEmpty(errno syscall.Errno) bool {
	return false
}
