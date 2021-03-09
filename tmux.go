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

type TmuxWindow struct {
	Name   string
	Layout string
	Root   string
}

type TmuxPane struct {
	Root string
	Type string
}

func (tmux Tmux) NewSession(name string, root string, windowName string) (string, error) {
	cmd := exec.Command("tmux", "new", "-Pd", "-s", name, "-n", windowName, "-c", root)
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
	cmd := exec.Command("tmux", "neww", "-Pd", "-t", target, "-c", root, "-F", "#{window_id}", "-n", name)

	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SendKeys(target string, command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", target, command, "Enter")
	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) Attach(target string, stdin *os.File, stdout *os.File, stderr *os.File) error {
	cmd := exec.Command("tmux", "attach", "-d", "-t", target)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) RenumberWindows(target string) error {
	cmd := exec.Command("tmux", "move-window", "-r", "-s", target, "-t", target)
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) SplitWindow(target string, splitType string, root string) (string, error) {
	args := []string{"split-window", "-Pd"}

	switch splitType {
	case VSplit:
		args = append(args, "-v")
	case HSplit:
		args = append(args, "-h")
	}

	args = append(args, []string{"-t", target, "-c", root, "-F", "#{pane_id}"}...)

	cmd := exec.Command("tmux", args...)

	pane, err := tmux.commander.Exec(cmd)
	if err != nil {
		return "", err
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

func (tmux Tmux) SwitchClient(target string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", target)
	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) ListWindows(target string) ([]TmuxWindow, error) {
	var windows []TmuxWindow

	cmd := exec.Command("tmux", "list-windows", "-F", "#{window_name};#{window_layout};#{pane_current_path}", "-t", target)
	out, err := tmux.commander.Exec(cmd)
	if err != nil {
		return windows, err
	}

	windowsList := strings.Split(out, "\n")

	for _, w := range windowsList {
		windowInfo := strings.Split(w, ";")
		window := TmuxWindow{
			Name:   windowInfo[0],
			Layout: windowInfo[1],
			Root:   windowInfo[2],
		}
		windows = append(windows, window)
	}

	return windows, nil

}

func (tmux Tmux) ListPanes(target string) ([]TmuxPane, error) {
	var panes []TmuxPane

	cmd := exec.Command("tmux", "list-panes", "-F", "", "-t", target)

	out, err := tmux.commander.Exec(cmd)
	if err != nil {
		return panes, err
	}

	panesList := strings.Split(out, "\n")

	for _, p := range panesList {
		paneInfo := strings.Split(p, ";")
		pane := TmuxPane{
			Type: paneInfo[0],
			Root: paneInfo[1],
		}

		panes = append(panes, pane)
	}

	return panes, nil
}
