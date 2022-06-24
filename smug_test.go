package main

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var testTable = map[string]struct {
	config           *Config
	options          *Options
	context          Context
	startCommands    []string
	stopCommands     []string
	commanderOutputs []string
}{
	"test with 1 window": {
		&Config{
			Session:     "ses",
			Root:        "~/root",
			BeforeStart: []string{"command1", "command2"},
			Windows: []Window{
				{
					Name:     "win1",
					Commands: []string{"command1"},
				},
			},
		},
		&Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"/bin/sh -c command1",
			"/bin/sh -c command2",
			"tmux new -Pd -s ses -n smug_def -c smug/root",
			"tmux neww -Pd -t ses: -c smug/root -F #{window_id} -n win1",
			"tmux send-keys -t win1 command1 Enter",
			"tmux select-layout -t win1 even-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{"ses", "win1"},
	},
	"test with 1 window and Detach: true": {
		&Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
			Windows: []Window{
				{
					Name: "win1",
				},
			},
		},
		&Options{
			Detach: true,
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
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{"xyz"},
	},
	"test with multiple windows and panes": {
		&Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{
					Name:   "win1",
					Manual: false,
					Layout: "main-horizontal",
					Panes: []Pane{
						{
							Type:     "horizontal",
							Commands: []string{"command1"},
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
		&Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n smug_def -c root",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux split-window -Pd -h -t win1 -c root -F #{pane_id}",
			"tmux send-keys -t win1.1 command1 Enter",
			"tmux select-layout -t win1 main-horizontal",
			"tmux kill-window -t ses:smug_def",
			"tmux move-window -r -s ses: -t ses:",
			"tmux attach -d -t ses:win1",
		},
		[]string{
			"/bin/sh -c stop1",
			"/bin/sh -c stop2 -d --foo=bar",
			"tmux kill-session -t ses",
		},
		[]string{"ses", "ses", "win1", "1"},
	},
	"test start windows from option's Windows parameter": {
		&Config{
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
		&Options{
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
		[]string{"xyz"},
	},
	"test attach to the existing session": {
		&Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{Name: "win1"},
			},
		},
		&Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux attach -d -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{""},
	},
	"test start a new session from another tmux session": {
		&Config{
			Session: "ses",
			Root:    "root",
		},
		&Options{Attach: false},
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
		[]string{"xyz"},
	},
	"test switch a client from another tmux session": {
		&Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{Name: "win1"},
			},
		},
		&Options{Attach: true},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux has-session -t ses:",
			"tmux switch-client -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{""},
	},
	"test create new windows in current session with same name": {
		&Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{Name: "win1"},
			},
		},
		&Options{
			InsideCurrentSession: true,
		},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux display-message -p #S",
			"tmux has-session -t ses:",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux select-layout -t  even-horizontal",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{"ses", ""},
	},
	"test create new windows in current session with different name": {
		&Config{
			Session: "ses",
			Root:    "root",
			Windows: []Window{
				{Name: "win1"},
			},
		},
		&Options{
			InsideCurrentSession: true,
		},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux display-message -p #S",
			"tmux has-session -t ses:",
			"tmux neww -Pd -t ses: -c root -F #{window_id} -n win1",
			"tmux select-layout -t win1 even-horizontal",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		[]string{"ses", "win1"},
	},
}

type MockCommander struct {
	Commands []string
	Outputs  []string
}

func (c *MockCommander) Exec(cmd *exec.Cmd) (string, error) {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))

	output := ""
	if len(c.Outputs) > 1 {
		output, c.Outputs = c.Outputs[0], c.Outputs[1:]
	} else if len(c.Outputs) == 1 {
		output = c.Outputs[0]
	}

	return output, nil
}

func (c *MockCommander) ExecSilently(cmd *exec.Cmd) error {
	c.Commands = append(c.Commands, strings.Join(cmd.Args, " "))
	return nil
}

func TestStartStopSession(t *testing.T) {
	os.Setenv("HOME", "smug") // Needed for testing ExpandPath function

	for testDescription, params := range testTable {

		t.Run("start session: "+testDescription, func(t *testing.T) {
			commander := &MockCommander{[]string{}, params.commanderOutputs}
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

		t.Run("stop session: "+testDescription, func(t *testing.T) {
			commander := &MockCommander{[]string{}, params.commanderOutputs}
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

func TestPrintCurrentSession(t *testing.T) {
	expectedConfig := Config{
		Session: "session_name",
		Windows: []Window{
			{
				Name:   "win1",
				Root:   "root",
				Layout: "layout",
				Panes: []Pane{
					{},
					{
						Root: "/tmp",
					},
				},
			},
		},
	}

	commander := &MockCommander{[]string{}, []string{
		"session_name",
		"id1;win1;layout;root",
		"root\n/tmp",
	}}
	tmux := Tmux{commander}

	smug := Smug{tmux, commander}

	actualConfig, err := smug.GetConfigFromSession(&Options{Project: "test"}, Context{})
	if err != nil {
		t.Fatalf("error %v", err)
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf("expected %v, got %v", expectedConfig, actualConfig)
	}
}
