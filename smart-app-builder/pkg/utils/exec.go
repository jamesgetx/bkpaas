package utils

import (
	"os"
	"os/exec"
)

// RunCommand runs a command
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
