package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const version = "v0.1.8"

var usage = fmt.Sprintf(`Smug - tmux session manager. Version %s


Usage:
	smug <command> [<project>] [-f, --file <file>] [-w, --windows <window>]... [-a, --attach] [-d, --debug]

Options:
	-f, --file %s
	-w, --windows %s
	-a, --attach %s
	-d, --debug %s

Examples:
	$ smug start blog
	$ smug start blog:win1
	$ smug start blog -w win1
	$ smug start blog:win1,win2
	$ smug stop blog
	$ smug start blog --attach
`, version, FileUsage, WindowsUsage, AttachUsage, DebugUsage)

func main() {
	options, err := ParseOptions(os.Args[1:], func() {
		fmt.Fprintf(os.Stdout, usage)
		os.Exit(0)
	})

	if err == ErrHelp {
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse command line options: %q", err.Error())
		os.Exit(1)
	}

	userConfigDir := filepath.Join(ExpandPath("~/"), ".config/smug")

	var configPath string
	if options.Config != "" {
		configPath = options.Config
	} else {
		configPath = filepath.Join(userConfigDir, options.Project+".yml")
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
	case CommandStart:
		if len(options.Windows) == 0 {
			fmt.Println("Starting a new session...")
		} else {
			fmt.Println("Starting new windows...")
		}
		config, err := GetConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = smug.Start(config, options, context)
		if err != nil {
			fmt.Println("Oops, an error occurred! Rolling back...")
			smug.Stop(config, options, context)
			os.Exit(1)
		}
	case CommandStop:
		if len(options.Windows) == 0 {
			fmt.Println("Terminating session...")
		} else {
			fmt.Println("Killing windows...")
		}
		config, err := GetConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = smug.Stop(config, options, context)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	case CommandNew, CommandEdit:
		err := EditConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

}
