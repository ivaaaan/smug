package main

import (
	"errors"
	"reflect"
	"testing"
)

var usageTestTable = []struct {
	argv []string
	opts Options
	err  error
}{
	{
		[]string{"start", "smug"},
		Options{
			Command:  "start",
			Project:  "smug",
			Config:   "",
			Windows:  []string{},
			Attach:   false,
			Detach:   false,
			Debug:    false,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "smug", "-w", "foo"},
		Options{
			Command:  "start",
			Project:  "smug",
			Config:   "",
			Windows:  []string{"foo"},
			Attach:   false,
			Detach:   false,
			Debug:    false,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "smug:foo,bar"},
		Options{
			Command:  "start",
			Project:  "smug",
			Config:   "",
			Windows:  []string{"foo", "bar"},
			Attach:   false,
			Detach:   false,
			Debug:    false,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "smug", "--attach", "--debug", "--detach"},
		Options{
			Command:  "start",
			Project:  "smug",
			Config:   "",
			Windows:  []string{},
			Attach:   true,
			Detach:   true,
			Debug:    true,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "smug", "-ad"},
		Options{
			Command:  "start",
			Project:  "smug",
			Config:   "",
			Windows:  []string{},
			Attach:   true,
			Detach:   false,
			Debug:    true,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "-f", "test.yml"},
		Options{
			Command:  "start",
			Project:  "",
			Config:   "test.yml",
			Windows:  []string{},
			Attach:   false,
			Detach:   false,
			Debug:    false,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "-f", "test.yml", "-w", "win1", "-w", "win2"},
		Options{
			Command:  "start",
			Project:  "",
			Config:   "test.yml",
			Windows:  []string{"win1", "win2"},
			Attach:   false,
			Detach:   false,
			Debug:    false,
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"start", "project", "a=b", "x=y"},
		Options{
			Command: "start",
			Project: "project",
			Config:  "",
			Windows: []string{},
			Attach:  false,
			Detach:  false,
			Debug:   false,
			Settings: map[string]string{
				"a": "b",
				"x": "y",
			},
		},
		nil,
	},
	{
		[]string{"start", "-f", "test.yml", "a=b", "x=y"},
		Options{
			Command: "start",
			Project: "",
			Config:  "test.yml",
			Windows: []string{},
			Attach:  false,
			Detach:  false,
			Debug:   false,
			Settings: map[string]string{
				"a": "b",
				"x": "y",
			},
		},
		nil,
	},
	{
		[]string{"start", "-f", "test.yml", "-w", "win1", "-w", "win2", "a=b", "x=y"},
		Options{
			Command: "start",
			Project: "",
			Config:  "test.yml",
			Windows: []string{"win1", "win2"},
			Attach:  false,
			Detach:  false,
			Debug:   false,
			Settings: map[string]string{
				"a": "b",
				"x": "y",
			},
		},
		nil,
	},
	{
		[]string{"start", "--help"},
		Options{},
		ErrHelp,
	},
	{
		[]string{"test"},
		Options{
			Command:  "start",
			Project:  "test",
			Windows:  []string{},
			Settings: map[string]string{},
		},
		nil,
	},
	{
		[]string{"test", "-w", "win1", "-w", "win2", "a=b", "x=y"},
		Options{
			Command:  "start",
			Project:  "test",
			Windows:  []string{"win1", "win2"},
			Settings: map[string]string{"a": "b", "x": "y"},
		},
		nil,
	},
	{
		[]string{},
		Options{},
		ErrHelp,
	},
	{
		[]string{"--help"},
		Options{},
		ErrHelp,
	},
	{
		[]string{"start", "--test"},
		Options{},
		errors.New("unknown flag: --test"),
	},
}

func TestParseOptions(t *testing.T) {
	for _, v := range usageTestTable {
		opts, err := ParseOptions(v.argv)
		if v.err != nil && err != nil && err.Error() != v.err.Error() {
			t.Errorf("expected error %v, got %v", v.err, err)
		}

		if opts != nil && !reflect.DeepEqual(v.opts, *opts) {
			t.Errorf("expected struct %v, got %v", v.opts, opts)
		}
	}
}
