package main

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var startSessionTestTable = []struct {
	config   Config
	commands []string
	windows  []string
}{
	{
		Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
		},
		[]string{
			"tmux has-session -t ses",
			"command1",
			"command2",
			"tmux new -Pd -s ses",
			"tmux kill-window -t ses:0",
			"tmux move-window -r",
			"tmux attach -t ses:0",
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
			"tmux new -Pd -s ses",
			"tmux neww -Pd -t ses: -n win1 -c root",
			"tmux kill-window -t ses:0",
			"tmux move-window -r",
			"tmux attach -t ses:0",
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
			"tmux new -Pd -s ses",
			"tmux neww -Pd -t ses: -n win2 -c root",
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
	for _, params := range startSessionTestTable {
		commander := &MockCommander{}
		tmux := Tmux{commander}
		smug := Smug{tmux, commander}

		smug.StartSession(params.config, params.windows)

		if !reflect.DeepEqual(params.commands, commander.Commands) {
			t.Errorf("expected\n%s\ngot\n%s", strings.Join(params.commands, "\n"), strings.Join(commander.Commands, "\n"))
		}
	}
}
