package main

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var testTable = []struct {
	config          Config
	options         Options
	context         Context
	startCommands   []string
	stopCommands    []string
	commanderOutput string
}{
	{
		Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
			Windows: []Window{
				{
					Name: "win1",
				},
			},
		},
		Options{
			Windows: []string{},
		},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"/bin/sh -c command1",
			"/bin/sh -c command2",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux select-layout -t xyz even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"xyz",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name:   "win1",
					Manual: false,
					Layout: "main-horizontal",
					Panes: []Pane{
						{
							Type: "horizontal",
						},
					},
				},
				{
					Name:   "win2",
					Manual: true,
					Layout: "tiled",
				},
			},
			Stop: []string{
				"stop1",
				"stop2 -d --foo=bar",
			},
		},
		Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux split-window -Pd -h -t 1 -c root -F #{pane_id}",
			"tmux select-layout -t 1 main-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"/bin/sh -c stop1",
			"/bin/sh -c stop2 -d --foo=bar",
			"tmux kill-session -t ses",
		},
		"1",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name:   "win1",
					Manual: false,
				},
				{
					Name:   "win2",
					Manual: true,
				},
			},
		},
		Options{
			Windows: []string{"win2"},
		},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win2",
			"tmux select-layout -t xyz even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
		},
		[]string{
			"tmux kill-window -t ses:win2",
		},
		"xyz",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name:     "win1",
					Manual:   false,
					Commands: []string{"command1", "command2"},
				},
				{
					Name:     "win2",
					Manual:   false,
					Commands: []string{"command3", "command4"},
				},
			},
		},
		Options{
			Windows: []string{},
		},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux send-keys -t xyz command1 Enter",
			"tmux send-keys -t xyz command2 Enter",
			"tmux select-layout -t xyz even-horizontal",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win2",
			"tmux send-keys -t xyz command3 Enter",
			"tmux send-keys -t xyz command4 Enter",
			"tmux select-layout -t xyz even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"xyz",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name:   "win1",
					Manual: false,
					Root:   "./win1",
					Panes: []Pane{
						{
							Root: "pane1",
							Type: "vertical",
							Commands: []string{
								"command1",
							},
						},
					},
				},
			},
		},
		Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root/win1 -F #{window_id} -n win1",
			"tmux split-window -Pd -v -t 1 -c root/win1/pane1 -F #{pane_id}",
			"tmux send-keys -t 1.1 command1 Enter",
			"tmux select-layout -t 1 even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"1",
	},
	{
		Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
			Windows: []Window{
				{Name: "win1"},
			},
		},
		Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux attach -d -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name: "win1",
				},
			},
		},
		Options{Attach: true},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux select-layout -t xyz even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux switch-client -t ses:win1",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"xyz",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
		},
		Options{Attach: false},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"xyz",
	},
	{
		Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{Name: "win1"},
			},
		},
		Options{Attach: true},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux has-session -t ses:",
			"tmux switch-client -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"",
	},
}

type MockCommander struct {
	Commands      []string
	DefaultOutput string
}

func (c *MockCommander) Exec(cmd *exec.Cmd) (string, error) {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))

	return c.DefaultOutput, nil
}

func (c *MockCommander) ExecSilently(cmd *exec.Cmd) error {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))
	return nil
}

func TestStartSession(t *testing.T) {
	for _, params := range testTable {

		t.Run("test start session", func(t *testing.T) {
			commander := &MockCommander{[]string{}, params.commanderOutput}
			tmux := Tmux{commander}
			smug := Smug{tmux, commander}

			err := smug.Start(params.config, params.options, params.context)
			if err != nil {
				t.Fatalf("error %v", err)
			}

			if !reflect.DeepEqual(params.startCommands, commander.Commands) {
				t.Errorf("expected\n%s\ngot\n%s", strings.Join(params.startCommands, "\n"), strings.Join(commander.Commands, "\n"))
			}
		})

		t.Run("test stop session", func(t *testing.T) {
			commander := &MockCommander{[]string{}, params.commanderOutput}
			tmux := Tmux{commander}
			smug := Smug{tmux, commander}

			err := smug.Stop(params.config, params.options, params.context)
			if err != nil {
				t.Fatalf("error %v", err)
			}

			if !reflect.DeepEqual(params.stopCommands, commander.Commands) {
				t.Errorf("expected\n%s\ngot\n%s", strings.Join(params.stopCommands, "\n"), strings.Join(commander.Commands, "\n"))
			}
		})

	}
}
