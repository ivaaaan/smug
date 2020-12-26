package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	var logger *log.Logger
	if options.Debug {
		logFile, err := os.Create(filepath.Join(userConfigDir, "smug.log"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		logger = log.New(logFile, "", 0)
	}

	commander := DefaultCommander{logger}
	tmux := Tmux{commander}
	smug := Smug{tmux, commander}

	context := CreateContext()

	switch options.Command {
	case "start":
		if len(options.Windows) == 0 {
			fmt.Println("Starting a new session...")
		} else {
			fmt.Println("Starting new windows...")
		}
		err = smug.Start(*config, options, *context)
		if err != nil {
			fmt.Println("Oops, an error occurred! Rolling back...")
			smug.Stop(*config, options, *context)
		}
	case "stop":
		if len(options.Windows) == 0 {
			fmt.Println("Terminating session...")
		} else {
			fmt.Println("Killing windows...")
		}
		err = smug.Stop(*config, options, *context)
	default:
		err = fmt.Errorf("Unknown command %q", options.Command)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
