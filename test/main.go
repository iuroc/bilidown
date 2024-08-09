package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ffmpeg")
	bs, _ := cmd.CombinedOutput()
	fmt.Println(string(bs))
}
