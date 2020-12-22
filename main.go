package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

func main() {
	parser := docopt.Parser{}

	options, err := ParseOptions(parser, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse command line otions: %q", err.Error())
		os.Exit(1)
	}

	userConfigDir := filepath.Join(ExpandPath("~/"), ".config/smug")
	configPath := filepath.Join(userConfigDir, options.Project+".yml")

	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	config, err := ParseConfig(string(f))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	commander := DefaultCommander{}
	tmux := Tmux{commander}
	smug := Smug{tmux, commander}

	switch options.Command {
	case "start":
		fmt.Println("Starting a new session...")
		err = smug.Start(*config, options.Windows)
		if err != nil {
			fmt.Println("Oops, an error occurred! Rolling back...")
			smug.Stop(*config, options.Windows)
		}
	case "stop":
		if len(options.Windows) == 0 {
			fmt.Println("Terminating session...")
		} else {
			fmt.Println("Killing windows...")
		}
		err = smug.Stop(*config, options.Windows)
	default:
		err = fmt.Errorf("Unknown command %q", options.Command)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
