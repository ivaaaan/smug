package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return strings.Replace(path, "~", userHome, 1)
	}

	return path
}

func Contains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}

	return false
}

type Smug struct {
	tmux       Tmux
	commander  Commander
	configPath string
}

func (smug Smug) execShellCommands(commands []string, path string) error {
	for _, c := range commands {
		cmd := exec.Command("/bin/sh", "-c", c)
		cmd.Dir = path

		_, err := smug.commander.Exec(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (smug Smug) switchOrAttach(sessionName string, attach bool, insideTmuxSession bool) error {
	if insideTmuxSession && attach {
		return smug.tmux.SwitchClient(sessionName)
	} else if !insideTmuxSession {
		return smug.tmux.Attach(sessionName, os.Stdin, os.Stdout, os.Stderr)
	}
	return nil
}

func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (smug Smug) Create() error {
	exists, err := IsFileExists(smug.configPath)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("File already exists")
	}
	file, err := os.Create(smug.configPath)
	defer file.Close()
	return err
}

func (smug Smug) Edit() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, smug.configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func (smug Smug) Stop(config Config, options Options, context Context) error {
	windows := options.Windows
	if len(windows) == 0 {
		sessionRoot := ExpandPath(config.Root)

		err := smug.execShellCommands(config.Stop, sessionRoot)
		if err != nil {
			return err
		}
		_, err = smug.tmux.StopSession(config.Session)
		return err
	}

	for _, w := range windows {
		err := smug.tmux.KillWindow(config.Session + ":" + w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (smug Smug) Start(config Config, options Options, context Context) error {
	sessionName := config.Session + ":"
	sessionExists := smug.tmux.SessionExists(sessionName)
	sessionRoot := ExpandPath(config.Root)

	windows := options.Windows
	attach := options.Attach

	if !sessionExists {
		err := smug.execShellCommands(config.BeforeStart, sessionRoot)
		if err != nil {
			return err
		}

		var defaultWindowName string
		if len(windows) > 0 {
			defaultWindowName = windows[0]
		} else if len(config.Windows) > 0 {
			defaultWindowName = config.Windows[0].Name
		}

		_, err = smug.tmux.NewSession(strings.Replace(sessionName, ":", "", 1), sessionRoot, defaultWindowName)
		if err != nil {
			return err
		}
	} else if len(windows) == 0 {
		return smug.switchOrAttach(sessionName, attach, context.InsideTmuxSession)
	}

	for wIndex, w := range config.Windows {
		if (len(windows) == 0 && w.Manual) || (len(windows) > 0 && !Contains(windows, w.Name)) {
			continue
		}

		windowRoot := ExpandPath(w.Root)
		if windowRoot == "" || !filepath.IsAbs(windowRoot) {
			windowRoot = filepath.Join(sessionRoot, w.Root)
		}

		window := sessionName + w.Name
		if (!sessionExists && wIndex > 0 && len(windows) == 0) || (sessionExists && len(windows) > 0) {
			_, err := smug.tmux.NewWindow(sessionName, w.Name, windowRoot)
			if err != nil {
				return err
			}
		}

		for _, c := range w.Commands {
			err := smug.tmux.SendKeys(window, c)
			if err != nil {
				return err
			}
		}

		for _, p := range w.Panes {
			paneRoot := ExpandPath(p.Root)
			if paneRoot == "" || !filepath.IsAbs(p.Root) {
				paneRoot = filepath.Join(windowRoot, p.Root)
			}

			_, err := smug.tmux.SplitWindow(window, p.Type, paneRoot, p.Commands)
			if err != nil {
				return err
			}
		}

		layout := w.Layout
		if layout == "" {
			layout = EvenHorizontal
		}

		_, err := smug.tmux.SelectLayout(sessionName+w.Name, layout)
		if err != nil {
			return err
		}
	}

	if len(windows) == 0 {
		return smug.switchOrAttach(sessionName, attach, context.InsideTmuxSession)
	}

	return nil
}
