package main

import (
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
		Options{"start", "smug", []string{}, false, false},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "-w", "foo"},
		Options{"start", "smug", []string{"foo"}, false, false},
		nil,
		0,
	},
	{
		[]string{"start", "smug:foo,bar"},
		Options{"start", "smug", []string{"foo", "bar"}, false, false},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "--attach", "--debug"},
		Options{"start", "smug", []string{}, true, true},
		nil,
		0,
	},
	{
		[]string{"start", "smug", "-ad"},
		Options{"start", "smug", []string{}, true, true},
		nil,
		0,
	},
	{
		[]string{"start"},
		Options{},
		ErrHelp,
		1,
	},
	{
		[]string{"start", "--help"},
		Options{},
		ErrHelp,
		1,
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

		if err != v.err {
			t.Errorf("expected to get error %v, got %v", v.err, err)
		}

		helpCalls = 0
	}
}
