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
	Tiled          = "tiled"
)

type TmuxOptions struct {
	// Default socket name
	SocketName string `yaml:"socket_name"`

	// Default socket path, overrides SocketName
	SocketPath string `yaml:"socket_path"`

	// tmux config file
	ConfigFile string `yaml:"config_file"`
}

type Tmux struct {
	commander Commander
	*TmuxOptions
}

type TmuxWindow struct {
	ID     string
	Name   string
	Layout string
	Root   string
}

type TmuxPane struct {
	Root string
}

func (tmux Tmux) cmd(args ...string) *exec.Cmd {
	tmuxCmd := []string{"tmux"}
	if tmux.SocketPath != "" {
		tmuxCmd = append(tmuxCmd, "-S", tmux.SocketPath)
	} else if tmux.SocketName != "" {
		tmuxCmd = append(tmuxCmd, "-L", tmux.SocketName)
	}

	if tmux.ConfigFile != "" {
		tmuxCmd = append(tmuxCmd, "-f", tmux.ConfigFile)
	}

	tmuxCmd = append(tmuxCmd, args...)

	return exec.Command(tmuxCmd[0], tmuxCmd[1:]...)
}

func (tmux Tmux) NewSession(name string, root string, windowName string) (string, error) {
	cmd := tmux.cmd("new", "-Pd", "-s", name, "-n", windowName, "-c", root)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SessionExists(name string) bool {
	cmd := tmux.cmd("has-session", "-t", name)
	res, err := tmux.commander.Exec(cmd)
	return res == "" && err == nil
}

func (tmux Tmux) KillWindow(target string) error {
	cmd := tmux.cmd("kill-window", "-t", target)
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) SelectWindow(target string) error {
	cmd := tmux.cmd("select-window", "-t", target)
	_, err := tmux.commander.Exec(cmd)
	return err
}

func (tmux Tmux) NewWindow(target string, name string, root string) (string, error) {
	cmd := tmux.cmd("neww", "-Pd", "-t", target, "-c", root, "-F", "#{window_id}", "-n", name)

	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SendKeys(target string, command string) error {
	cmd := tmux.cmd("send-keys", "-t", target, command, "Enter")
	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) Attach(target string, stdin *os.File, stdout *os.File, stderr *os.File) error {
	cmd := tmux.cmd("attach", "-d", "-t", target)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) RenumberWindows(target string) error {
	cmd := tmux.cmd("move-window", "-r", "-s", target, "-t", target)
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

	cmd := tmux.cmd(args...)

	pane, err := tmux.commander.Exec(cmd)
	if err != nil {
		return "", err
	}

	return pane, nil
}

func (tmux Tmux) SelectLayout(target string, layoutType string) (string, error) {
	cmd := tmux.cmd("select-layout", "-t", target, layoutType)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SetEnv(target string, key string, value string) (string, error) {
	cmd := tmux.cmd("setenv", "-t", target, key, value)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) StopSession(target string) (string, error) {
	cmd := tmux.cmd("kill-session", "-t", target)
	return tmux.commander.Exec(cmd)
}

func (tmux Tmux) SwitchClient(target string) error {
	cmd := tmux.cmd("switch-client", "-t", target)
	return tmux.commander.ExecSilently(cmd)
}

func (tmux Tmux) SessionName() (string, error) {
	cmd := tmux.cmd("display-message", "-p", "#S")
	sessionName, err := tmux.commander.Exec(cmd)
	if err != nil {
		return sessionName, err
	}

	return sessionName, nil
}

func (tmux Tmux) ListWindows(target string) ([]TmuxWindow, error) {
	var windows []TmuxWindow

	cmd := tmux.cmd("list-windows", "-F", "#{window_id};#{window_name};#{window_layout};#{pane_current_path}", "-t", target)
	out, err := tmux.commander.Exec(cmd)
	if err != nil {
		return windows, err
	}

	windowsList := strings.Split(out, "\n")

	for _, w := range windowsList {
		windowInfo := strings.Split(w, ";")
		window := TmuxWindow{
			ID:     windowInfo[0],
			Name:   windowInfo[1],
			Layout: windowInfo[2],
			Root:   windowInfo[3],
		}
		windows = append(windows, window)
	}

	return windows, nil
}

func (tmux Tmux) ListPanes(target string) ([]TmuxPane, error) {
	var panes []TmuxPane

	cmd := tmux.cmd("list-panes", "-F", "#{pane_current_path}", "-t", target)

	out, err := tmux.commander.Exec(cmd)
	if err != nil {
		return panes, err
	}

	panesList := strings.Split(out, "\n")

	for _, p := range panesList {
		paneInfo := strings.Split(p, ";")
		pane := TmuxPane{
			Root: paneInfo[0],
		}

		panes = append(panes, pane)
	}

	return panes, nil
}
