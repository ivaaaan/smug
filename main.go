package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var version = "[dev build]"

var usage = fmt.Sprintf(`Smug - tmux session manager. Version %s


Usage:
	smug <command> [<project>] [-f, --file <file>] [-w, --windows <window>]... [-a, --attach] [-d, --debug] [--detach] [-i, --inside-current-session] [<key>=<value>]...

Options:
	-f, --file %s
	-w, --windows %s
	-a, --attach %s
	-i, --inside-current-session %s
	-d, --debug %s
	--detach %s

Commands:
	list    list available project configurations
	edit    edit project configuration
	new     new project configuration
	start   start project session
	stop    stop project session
	print   session configuration to stdout

Examples:
	$ smug list
	$ smug edit blog
	$ smug new blog
	$ smug start blog
	$ smug start blog:win1
	$ smug start blog -w win1
	$ smug start blog:win1,win2
	$ smug stop blog
	$ smug start blog --attach
	$ smug print > ~/.config/smug/blog.yml
`, version, FileUsage, WindowsUsage, AttachUsage, InsideCurrentSessionUsage, DebugUsage, DetachUsage)

const defaultConfigFile = ".smug.yml"

func newLogger(path string) *log.Logger {
	logFile, err := os.Create(filepath.Join(path, "smug.log"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return log.New(logFile, "", 0)
}

func main() {
	options, err := ParseOptions(os.Args[1:])
	if errors.Is(err, ErrHelp) {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Cannot parse command line options: %q",
			err.Error(),
		)
		os.Exit(1)
	}

	userConfigDir := filepath.Join(ExpandPath("~/"), ".config/smug")

	var logger *log.Logger
	if options.Debug {
		logger = newLogger(userConfigDir)
	}

	commander := DefaultCommander{logger}
	tmux := Tmux{commander, &TmuxOptions{}}
	smug := Smug{tmux, commander}
	context := CreateContext()

	var configPath string
	if options.Config != "" {
		configPath = options.Config
	} else if options.Project != "" {

		config, err := FindConfig(userConfigDir, options.Project)

		if err != nil && options.Command != CommandNew && options.Command != CommandStart && options.Command != CommandStop {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if options.Command == CommandNew {
			config = fmt.Sprintf("%s.yml", options.Project)
		}
		configPath = filepath.Join(userConfigDir, config)
	} else {
		path, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		configPath = filepath.Join(path, defaultConfigFile)
	}

	switch options.Command {
	case CommandStart:
		if len(options.Windows) == 0 {
			fmt.Println("Starting a new session...")
		} else {
			fmt.Println("Starting new windows...")
		}
		configs, err := FindConfigs(userConfigDir, options.Project)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		for configIndex, configPath := range configs {
			config, err := GetConfig(configPath, options.Settings, smug.tmux.TmuxOptions)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
			options.Detach = configIndex != len(configs)-1
			err = smug.Start(config, options, context)
			if err != nil {
				fmt.Println("Oops, an error occurred! Rolling back...")
				smug.Stop(config, options, context)
				os.Exit(1)
			}
		}
	case CommandStop:
		if len(options.Windows) == 0 {
			fmt.Println("Terminating session...")
		} else {
			fmt.Println("Killing windows...")
		}
		configs, err := FindConfigs(userConfigDir, options.Project)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		for _, configPath := range configs {
			config, err := GetConfig(configPath, options.Settings, smug.tmux.TmuxOptions)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}

			err = smug.Stop(config, options, context)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	case CommandNew, CommandEdit:
		err := EditConfig(configPath)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	case CommandList:
		configs, err := ListConfigs(userConfigDir, true)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		for _, config := range configs {
			fileExt := path.Ext(config)
			fmt.Println(strings.TrimSuffix(config, fileExt))
			isDir, err := IsDirectory(userConfigDir+"/"+config)
			if err != nil {
				continue
			}
			if isDir {

				subConfigs, err := ListConfigs(userConfigDir+"/"+config, false)
				if err != nil {
					fmt.Fprint(os.Stderr, err.Error())
					os.Exit(1)
				}
				for _, subConfig := range subConfigs {
					fileExt := path.Ext(subConfig)
					fmt.Println("|--"+strings.TrimSuffix(subConfig, fileExt))
				}

			}

		}

	case CommandPrint:
		config, err := smug.GetConfigFromSession(options, context)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		d, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Println(string(d))
	}
}
