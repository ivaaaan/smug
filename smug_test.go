package main

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var testTable = []struct {
	config        Config
	startCommands []string
	stopCommands  []string
	windows       []string
}{
	{
		Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
		},
		[]string{
			"tmux has-session -t ses",
			"/bin/sh -c command1",
			"/bin/sh -c command2",
			"tmux new -Pd -s ses -n ",
			"tmux attach -d -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{},
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
						Pane{
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
		[]string{
			"tmux has-session -t ses",
			"tmux new -Pd -s ses -n win1",
			"tmux split-window -Pd -t ses:win1 -c root -h",
			"tmux select-layout -t ses:win1 main-horizontal",
			"tmux attach -d -t ses:",
		},
		[]string{
			"/bin/sh -c stop1",
			"/bin/sh -c stop2 -d --foo=bar",
			"tmux kill-session -t ses",
		},
		[]string{},
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
		[]string{
			"tmux has-session -t ses",
			"tmux new -Pd -s ses -n win2",
			"tmux select-layout -t ses:win2 even-horizontal",
		},
		[]string{
			"tmux kill-window -t ses:win2",
		},
		[]string{
			"win2",
		},
	},
}

type MockCommander struct {
	Commands []string
}

func (c *MockCommander) Exec(cmd *exec.Cmd) (string, error) {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))

	return "ses:", nil
}

func (c *MockCommander) ExecSilently(cmd *exec.Cmd) error {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))
	return nil
}

func TestStartSession(t *testing.T) {
	for _, params := range testTable {

		t.Run("test start session", func(t *testing.T) {
			commander := &MockCommander{}
			tmux := Tmux{commander}
			smug := Smug{tmux, commander}

			err := smug.Start(params.config, params.windows, false)
			if err != nil {
				t.Fatalf("error %v", err)
			}

			if !reflect.DeepEqual(params.startCommands, commander.Commands) {
				t.Errorf("expected\n%s\ngot\n%s", strings.Join(params.startCommands, "\n"), strings.Join(commander.Commands, "\n"))
			}
		})

		t.Run("test stop session", func(t *testing.T) {
			commander := &MockCommander{}
			tmux := Tmux{commander}
			smug := Smug{tmux, commander}

			err := smug.Stop(params.config, params.windows)
			if err != nil {
				t.Fatalf("error %v", err)
			}

			if !reflect.DeepEqual(params.stopCommands, commander.Commands) {
				t.Errorf("expected\n%s\ngot\n%s", strings.Join(params.stopCommands, "\n"), strings.Join(commander.Commands, "\n"))
			}
		})

	}
}
