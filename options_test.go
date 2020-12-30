package main

import (
	"reflect"
	"testing"
)

var usageTestTable = []struct {
	argv []string
	opts Options
}{
	{
		[]string{"start", "smug"},
		Options{"start", "smug", []string{}, false, false},
	},
	{
		[]string{"start", "smug", "-w", "foo"},
		Options{"start", "smug", []string{"foo"}, false, false},
	},
	{
		[]string{"start", "smug:foo,bar"},
		Options{"start", "smug", []string{"foo", "bar"}, false, false},
	},
	{
		[]string{"start", "smug", "--attach", "--debug"},
		Options{"start", "smug", []string{}, true, true},
	},
	{
		[]string{"start", "smug", "-ad"},
		Options{"start", "smug", []string{}, true, true},
	},
}

func TestParseOptions(t *testing.T) {
	for _, v := range usageTestTable {
		opts, err := ParseOptions(v.argv)

		if err != nil {
			t.Fail()
		}

		if !reflect.DeepEqual(v.opts, opts) {
			t.Errorf("expected struct %v, got %v", v.opts, opts)
		}
	}
}
