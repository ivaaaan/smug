package main

import (
	"errors"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

const (
	CommandStart = "start"
	CommandStop  = "stop"
	CommandNew   = "new"
	CommandEdit  = "edit"
	CommandList  = "list"
	CommandPrint = "print"
)

type command struct {
	Name    string
	Aliases []string
}

type commands []command

var Commands = commands{
	{
		Name:    CommandStart,
		Aliases: []string{},
	},
	{
		Name:    CommandStop,
		Aliases: []string{"s", "st"},
	},
	{
		Name:    CommandNew,
		Aliases: []string{"n"},
	},
	{
		Name:    CommandEdit,
		Aliases: []string{"e"},
	},
	{
		Name:    CommandList,
		Aliases: []string{"l"},
	},
	{
		Name:    CommandPrint,
		Aliases: []string{"p"},
	},
}

func (c *commands) Resolve(v string) (*command, error) {
	for _, cmd := range *c {
		if cmd.Name == v || Contains(cmd.Aliases, v) {
			return &cmd, nil
		}
	}

	return nil, ErrCommandNotFound
}

func (c *commands) FindByName(n string) *command {
	for _, cmd := range *c {
		if cmd.Name == n {
			return &cmd
		}
	}

	return nil
}

type Options struct {
	Command              string
	Project              string
	Config               string
	Windows              []string
	Settings             map[string]string
	Attach               bool
	Detach               bool
	Debug                bool
	InsideCurrentSession bool
}

var ErrHelp = errors.New("help requested")
var ErrCommandNotFound = errors.New("command not found")

const (
	WindowsUsage              = "List of windows to start. If session exists, those windows will be attached to current session"
	AttachUsage               = "Force switch client for a session"
	DetachUsage               = "Detach tmux session. The same as -d flag in the tmux"
	DebugUsage                = "Print all commands to ~/.config/smug/smug.log"
	FileUsage                 = "A custom path to a config file"
	InsideCurrentSessionUsage = "Create all windows inside current session"
)

// Creates a new FlagSet.
// Moved it to a variable to be able to override it in the tests.
var NewFlagSet = func(cmd string) *pflag.FlagSet {
	f := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	return f
}

func ParseOptions(argv []string, helpRequested func()) (Options, error) {
	if len(argv) == 0 {
		helpRequested()
		return Options{}, ErrHelp
	}

	if argv[0] == "--help" || argv[0] == "-h" {
		helpRequested()
		return Options{}, ErrHelp
	}

	cmd, cmdErr := Commands.Resolve(argv[0])
	if errors.Is(cmdErr, ErrCommandNotFound) {
		cmd = Commands.FindByName(CommandStart)
	}

	flags := NewFlagSet(cmd.Name)

	config := flags.StringP("file", "f", "", FileUsage)
	windows := flags.StringArrayP("windows", "w", []string{}, WindowsUsage)
	attach := flags.BoolP("attach", "a", false, AttachUsage)
	detach := flags.Bool("detach", false, DetachUsage)
	debug := flags.BoolP("debug", "d", false, DebugUsage)
	insideCurrentSession := flags.BoolP("inside-current-session", "i", false, InsideCurrentSessionUsage)

	err := flags.Parse(argv)
	if err == pflag.ErrHelp {
		return Options{}, ErrHelp
	}

	if err != nil {
		return Options{}, err
	}

	// If config file flag is not set, and env is, use the env
	if len(*config) == 0 && len(os.Getenv("SMUG_SESSION_CONFIG_PATH")) > 0 {
		*config = os.Getenv("SMUG_SESSION_CONFIG_PATH")
	}

	var project string
	if *config == "" {
		if errors.Is(cmdErr, ErrCommandNotFound) {
			project = argv[0]
		} else if len(argv) > 1 {
			project = argv[1]
		}
	}

	if strings.Contains(project, ":") {
		parts := strings.Split(project, ":")
		project = parts[0]
		wl := strings.Split(parts[1], ",")
		windows = &wl
	}

	settings := make(map[string]string)
	userSettings := flags.Args()[1:]
	if len(userSettings) > 0 {
		for _, kv := range userSettings {
			s := strings.Split(kv, "=")
			if len(s) < 2 {
				continue
			}
			settings[s[0]] = s[1]
		}
	}

	return Options{
		Project:              project,
		Config:               *config,
		Command:              cmd.Name,
		Settings:             settings,
		Windows:              *windows,
		Attach:               *attach,
		Detach:               *detach,
		Debug:                *debug,
		InsideCurrentSession: *insideCurrentSession,
	}, nil
}
