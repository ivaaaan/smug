package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultWindowName = "smug_def"

// Very wisely picked default value,
// after which panes will be rebalanced for each `split-window`
// Helps with "no space for new pane" error
const defaultRebalancePanesThreshold = 5

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
	tmux      Tmux
	commander Commander
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

func (smug Smug) setEnvVariables(target string, env map[string]string) error {
	for key, value := range env {
		_, err := smug.tmux.SetEnv(target, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (smug Smug) switchOrAttach(target string, attach bool, insideTmuxSession bool) error {
	if insideTmuxSession && attach {
		return smug.tmux.SwitchClient(target)
	} else if !insideTmuxSession {
		return smug.tmux.Attach(target, os.Stdin, os.Stdout, os.Stderr)
	}

	return nil
}

func (smug Smug) Stop(config Config, options *Options, context Context) error {
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

func (smug Smug) Start(config Config, options *Options, context Context) error {
	var sessionName string
	var err error

	createWindowsInsideCurrSession := options.InsideCurrentSession
	if createWindowsInsideCurrSession && !context.InsideTmuxSession {
		return errors.New("cannot use -i flag outside of a tmux session")
	}

	sessionName = config.Session
	if createWindowsInsideCurrSession {
		sessionName, err = smug.tmux.SessionName()
		if err != nil {
			return err
		}
	}
	sessionName = sessionName + ":"

	sessionExists := smug.tmux.SessionExists(sessionName)
	sessionRoot := ExpandPath(config.Root)

	rebalancePanesThreshold := config.RebalanceWindowsThreshold
	if rebalancePanesThreshold == 0 {
		rebalancePanesThreshold = defaultRebalancePanesThreshold
	}

	windows := options.Windows
	attach := options.Attach

	if !sessionExists && !createWindowsInsideCurrSession {
		err := smug.execShellCommands(config.BeforeStart, sessionRoot)
		if err != nil {
			return err
		}

		_, err = smug.tmux.NewSession(config.Session, sessionRoot, defaultWindowName)
		if err != nil {
			return err
		}

		err = smug.setEnvVariables(config.Session, config.Env)
		if err != nil {
			return err
		}
	} else if len(windows) == 0 && !createWindowsInsideCurrSession {
		return smug.switchOrAttach(sessionName, attach, context.InsideTmuxSession)
	}

	for _, w := range config.Windows {
		if (len(windows) == 0 && w.Manual) || (len(windows) > 0 && !Contains(windows, w.Name)) {
			continue
		}

		windowRoot := ExpandPath(w.Root)
		if windowRoot == "" || !filepath.IsAbs(windowRoot) {
			windowRoot = filepath.Join(sessionRoot, w.Root)
		}

		window, err := smug.tmux.NewWindow(sessionName, w.Name, windowRoot)
		if err != nil {
			return err
		}

		for _, c := range w.Commands {
			err := smug.tmux.SendKeys(window, c)
			if err != nil {
				return err
			}
		}

		for pIndex, p := range w.Panes {
			paneRoot := ExpandPath(p.Root)
			if paneRoot == "" || !filepath.IsAbs(p.Root) {
				paneRoot = filepath.Join(windowRoot, p.Root)
			}

			newPane, err := smug.tmux.SplitWindow(window, p.Type, paneRoot)

			if err != nil {
				return err
			}

			for _, c := range p.Commands {
				err = smug.tmux.SendKeys(window+"."+newPane, c)
				if err != nil {
					return err
				}
			}

			if pIndex+1 >= rebalancePanesThreshold {
				_, err = smug.tmux.SelectLayout(window, Tiled)
				if err != nil {
					return err
				}

			}
		}

		layout := w.Layout
		if layout == "" {
			layout = EvenHorizontal
		}

		_, err = smug.tmux.SelectLayout(window, layout)
		if err != nil {
			return err
		}
	}

	if !options.InsideCurrentSession {
		smug.tmux.KillWindow(sessionName + defaultWindowName)
		smug.tmux.RenumberWindows(sessionName)
	}

	if len(windows) == 0 && len(config.Windows) > 0 && options.Detach == false {
		return smug.switchOrAttach(sessionName+config.Windows[0].Name, attach, context.InsideTmuxSession)
	}

	return nil
}

func (smug Smug) GetConfigFromSession(options *Options, context Context) (Config, error) {
	config := Config{}

	tmuxSession, err := smug.tmux.SessionName()
	if err != nil {
		return Config{}, err
	}
	config.Session = tmuxSession

	tmuxWindows, err := smug.tmux.ListWindows(options.Project)
	if err != nil {
		return Config{}, err
	}

	for _, w := range tmuxWindows {
		tmuxPanes, err := smug.tmux.ListPanes(options.Project + ":" + w.Id)
		if err != nil {
			return Config{}, err
		}

		panes := []Pane{}
		for _, p := range tmuxPanes {
			root := p.Root
			if root == w.Root {
				root = ""
			}
			panes = append(panes, Pane{
				Root: root,
			})
		}

		config.Windows = append(config.Windows, Window{
			Name:   w.Name,
			Layout: w.Layout,
			Root:   w.Root,
			Panes:  panes,
		})
	}

	return config, nil
}
