package task_test

import (
	"bufio"
	"fmt"
	"os/exec"
	"testing"
)

func TestFFMPEG(t *testing.T) {
	cmd := exec.Command("ffmpeg", "-i", `E:\bilidown\27.video`, "-i", `E:\bilidown\27.audio`, "-c:v", "copy", "-c:a", "copy", "-progress", "pipe:1", `E:\bilidown\27.mp4`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(">>>", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}

	t.Log("done")
}

func TestNum(t *testing.T) {
	t.Log(int64(100999) / 1000)
}
