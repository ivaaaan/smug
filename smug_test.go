package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
		},
		Options{
			Windows: []string{},
		},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"/bin/sh -c command1",
			"/bin/sh -c command2",
			"tmux new -Pd -s ses -n  -c root",
			"tmux attach -d -t ses:",
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
		Options{},
		Context{},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n win1 -c root",
			"tmux split-window -Pd -t ses:win1 -c root -h",
			"tmux select-layout -t ses:win1 main-horizontal",
			"tmux attach -d -t ses:",
		},
		[]string{
			"/bin/sh -c stop1",
			"/bin/sh -c stop2 -d --foo=bar",
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
			"tmux new -Pd -s ses -n win2 -c root",
			"tmux select-layout -t ses:win2 even-horizontal",
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
			"tmux new -Pd -s ses -n win1 -c root",
			"tmux send-keys -t ses:win1 command1 Enter",
			"tmux send-keys -t ses:win1 command2 Enter",
			"tmux select-layout -t ses:win1 even-horizontal",
			"tmux neww -Pd -t ses: -n win2 -c root",
			"tmux send-keys -t ses:win2 command3 Enter",
			"tmux send-keys -t ses:win2 command4 Enter",
			"tmux select-layout -t ses:win2 even-horizontal",
			"tmux attach -d -t ses:",
		},
		[]string{
			"tmux kill-session -t ses",
		},
		"xyz",
	},

	{
		Config{
			Session:     "ses",
			Root:        "root",
			BeforeStart: []string{"command1", "command2"},
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
		},
		Options{Attach: true},
		Context{InsideTmuxSession: true},
		[]string{
			"tmux has-session -t ses:",
			"tmux new -Pd -s ses -n  -c root",
			"tmux switch-client -t ses:",
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
			"tmux new -Pd -s ses -n  -c root",
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

		userConfigDir := filepath.Join(ExpandPath("~/"), ".config/smug")
		configPath := filepath.Join(userConfigDir, params.options.Project+".yml")

		t.Run("test start session", func(t *testing.T) {
			commander := &MockCommander{[]string{}, params.commanderOutput}
			tmux := Tmux{commander}
			smug := Smug{tmux, commander, configPath}

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
			smug := Smug{tmux, commander, configPath}

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

func TestSmug_CreateEdit(t *testing.T) {
	tmpdir := t.TempDir()
	os.Setenv("EDITOR", "/usr/bin/vim")

	type fields struct {
		tmux       Tmux
		commander  Commander
		configPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				configPath: filepath.Join(tmpdir, "test1.yml"),
			},
		},
		{
			fields: fields{
				configPath: filepath.Join(tmpdir, "test2.yml"),
			},
		},
		{
			fields: fields{
				configPath: filepath.Join(tmpdir, "test3.yml"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smug := Smug{
				tmux:       tt.fields.tmux,
				commander:  tt.fields.commander,
				configPath: tt.fields.configPath,
			}
			if err := smug.Create(); (err != nil) != tt.wantErr {
				t.Errorf("Smug.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			// if err := smug.Edit(); (err != nil) != tt.wantErr {
			// 	t.Errorf("Smug.Edit() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}
