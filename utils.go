package bilidown

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenUrlByBrowser open opens the specified URL in the default browser of the user.
func OpenUrlByBrowser(url string) error {
	fmt.Println(url)
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		// "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
