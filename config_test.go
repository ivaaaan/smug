package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	yaml := `
session: ${session}
windows:
  - layout: tiled
    commands:
      - echo 1
    panes:
      - commands:
        - echo 2
        - echo ${HOME}
        type: horizontal`

	config, err := ParseConfig(yaml, map[string]string{
		"session": "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := Config{
		Session: "test",
		Windows: []Window{
			{
				Layout:   "tiled",
				Commands: []string{"echo 1"},
				Panes: []Pane{
					{
						Type:     "horizontal",
						Commands: []string{"echo 2", "echo " + os.Getenv("HOME")},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("expected %v, got %v", expected, config)
	}
}
