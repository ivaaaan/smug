package main

import (
	"os"
	"os/exec"
	"strings"
)

const (
	VSplit = "vertical"
	HSplit = "horizontal"
)

const (
	EvenHorizontal = "even-horizontal"
	EvenVertical   = "even-vertical"
	MainHorizontal = "main-horizontal"
	MainVertical   = "main-vertical"
	Tiled          = "tiled"
)

type Tmux struct {
	commander Commander
}

func (tmux Tmux) NewSession(name string, windowName string) (string, error) {
	cmd := exec.Command("tmux", "new", "-Pd", "-s", name, "-n", windowName)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SessionExists(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", name)
	res, err := tmux.commander.Exec(cmd)
	return res == "" && err == nil
}

func (tmux Tmux) KillWindow(target string) error {
	cmd := exec.Command("tmux", "kill-window", "-t", target)
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) NewWindow(target string, name string, root string) (string, error) {
	cmd := exec.Command("tmux", "neww", "-Pd", "-t", target, "-n", name, "-c", root)

	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SendKeys(target string, command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", target, command, "Enter")
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) Attach(target string, stdin *os.File, stdout *os.File, stderr *os.File) error {
	cmd := exec.Command("tmux", "attach", "-d", "-t", target)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) RenumberWindows(target string) error {
	cmd := exec.Command("tmux", "move-window", "-r")
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) SplitWindow(target string, splitType string, root string, commands []string) (string, error) {
	args := []string{"split-window", "-Pd", "-t", target, "-c", root}

	switch splitType {
	case VSplit:
		args = append(args, "-v")
	case HSplit:
		args = append(args, "-h")
	}

	cmd := exec.Command("tmux", args...)

	pane, err := tmux.commander.Exec(cmd)
	if err != nil {
		return "", err
	}

	for _, c := range commands {
		err = tmux.SendKeys(pane, c)
		if err != nil {
			return "", err
		}
	}

	return pane, nil
}

func (tmux Tmux) SelectLayout(target string, layoutType string) (string, error) {
	cmd := exec.Command("tmux", "select-layout", "-t", target, layoutType)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) StopSession(target string) (string, error) {
	cmd := exec.Command("tmux", "kill-session", "-t", target)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) ListWindows(target string) ([]string, error) {
	cmd := exec.Command("tmux", "list-windows", "-t", target, "-F", "#{window_index}")

	output, err := tmux.commander.Exec(cmd)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(output, "\n"), nil
}

func (tmux Tmux) SwitchClient(target string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", target)
	return tmux.commander.ExecSilently(cmd)
}
