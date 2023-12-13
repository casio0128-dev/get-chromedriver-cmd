package main

import (
	"runtime"
	"strings"
)

func getOSKey() string {
	switch runtime.GOOS {
	case "darwin":
		if strings.EqualFold(runtime.GOARCH, "386") {
			return "mac-x64"
		} else {
			return "mac-arm64"
		}
	case "windows":
		if strings.EqualFold(runtime.GOARCH, "386") {
			return "win32"
		} else {
			return "win64"
		}
	case "linux":
		fallthrough
	default:
		return "linux"
	}
	return ""
}
