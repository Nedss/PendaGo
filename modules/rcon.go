package modules

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunRCONCommand(host, port, password, command string) (string, error) {
	if host == "" || port == "" || password == "" {
		return "", fmt.Errorf("missing rcon configuration")
	}
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command is empty")
	}

	args := []string{
		"-H", host,
		"-P", port,
		"-p", password,
		command,
	}
	out, err := exec.Command("mcrcon", args...).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
