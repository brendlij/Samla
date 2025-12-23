package main

import (
	"os/exec"
	"runtime"
)

// openFolder attempts to show the given path in the platform file explorer.
func openFolder(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
