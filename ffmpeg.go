package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// FFMpeg sets up an ffmpeg process to convert stdin to opus on stdout
func FFMpeg(stdout io.Writer, stdin io.ReadCloser) error {
	defer func() {
		if err := stdin.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	cmd := exec.Command(
		"ffmpeg", "-hide_banner", "-loglevel", "warning", "-i", "-", "-f", "opus",
		"-vn", "-c:a", "libopus", "-b:a", "16k", "-application", "voip", "-",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = stdout
	cmd.Stdin = stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
