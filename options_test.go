package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

var usageTestTable = []struct {
	argv      []string
	opts      Options
	err       error
	helpCalls int
}{
	{
		[]string{"start", "smug"},
		Options{
			Command: "start",
			Project: "smug",
			Config:  "",
			Windows: []string{},
			Attach:  false,
			Detach:  false,
			Debug:   false,
		},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "-w", "foo"},
		Options{
			Command: "start",
			Project: "smug",
			Config:  "",
			Windows: []string{"foo"},
			Attach:  false,
			Detach:  false,
			Debug:   false,
		},
		nil,
		0,
	},
	{
		[]string{"start", "smug:foo,bar"},
		Options{
			Command: "start",
			Project: "smug",
			Config:  "",
			Windows: []string{"foo", "bar"},
			Attach:  false,
			Detach:  false,
			Debug:   false,
		},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "--attach", "--debug", "--detach"},
		Options{
			Command: "start",
			Project: "smug",
			Config:  "",
			Windows: []string{},
			Attach:  true,
			Detach:  true,
			Debug:   true,
		},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "-ad"},
		Options{
			Command: "start",
			Project: "smug",
			Config:  "",
			Windows: []string{},
			Attach:  true,
			Detach:  false,
			Debug:   true,
		},
		nil,
		0,
	},
	{
		[]string{"start", "-f", "test.yml"},
		Options{
			Command: "start",
			Project: "",
			Config:  "test.yml",
			Windows: []string{},
			Attach:  false,
			Detach:  false,
			Debug:   false,
		},
		nil,
		0,
	},
	{
		[]string{"start", "-f", "test.yml", "-w", "win1", "-w", "win2"},
		Options{
			Command: "start",
			Project: "",
			Config:  "test.yml",
			Windows: []string{"win1", "win2"},
			Attach:  false,
			Detach:  false,
			Debug:   false,
		},
		nil,
		0,
	},
	{
		[]string{"start", "--help"},
		Options{},
		ErrHelp,
		1,
	},
	{
		[]string{"test"},
		Options{},
		ErrHelp,
		1,
	},
	{
		[]string{},
		Options{},
		ErrHelp,
		1,
	},
	{
		[]string{"--help"},
		Options{},
		ErrHelp,
		1,
	},
	{
		[]string{"start", "--test"},
		Options{},
		errors.New("unknown flag: --test"),
		0,
	},
}

func TestParseOptions(t *testing.T) {
	helpCalls := 0
	helpRequested := func() {
		helpCalls++
	}

	NewFlagSet = func(cmd string) *pflag.FlagSet {
		flagSet := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
		flagSet.Usage = helpRequested
		return flagSet
	}

	for _, v := range usageTestTable {
		opts, err := ParseOptions(v.argv, helpRequested)

		if !reflect.DeepEqual(v.opts, opts) {
			t.Errorf("expected struct %v, got %v", v.opts, opts)
		}

		if helpCalls != v.helpCalls {
			t.Errorf("expected to get %d help calls, got %d", v.helpCalls, helpCalls)
		}

		if v.err != nil && err.Error() != v.err.Error() {
			t.Errorf("expected to get error %v, got %v", v.err, err)
		}

		helpCalls = 0
	}
}
