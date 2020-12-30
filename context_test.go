package main

import (
	"os"
	"reflect"
	"testing"
)

var environmentTestTable = []struct {
	environment map[string]string
	context     Context
}{
	{
		map[string]string{},
		Context{InsideTmuxSession: false},
	},
	{
		map[string]string{
			"TMUX": "",
		},
		Context{InsideTmuxSession: true},
	},
	{
		map[string]string{
			"TERM": "screen",
		},
		Context{InsideTmuxSession: true},
	},
	{
		map[string]string{
			"TERM": "xterm",
			"TMUX": "",
		},
		Context{InsideTmuxSession: true},
	},
}

func TestCreateContext(t *testing.T) {
	os.Clearenv()
	for _, v := range environmentTestTable {
		for key, value := range v.environment {
			os.Setenv(key, value)
		}

		context := CreateContext()

		if !reflect.DeepEqual(v.context, context) {
			t.Errorf("expected context %v, got %v", v.context, context)
		}
	}
}
