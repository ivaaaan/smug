package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var version = "[dev build]"

type Options struct {
	Project  string
	Config   string
	Windows  []string
	Settings map[string]string
	Attach   bool
	Detach   bool
	Debug    bool
}

func initOptions(c *cli.Context) Options {

	settings := make(map[string]string)

	userSettings := c.Args().Slice()[1:]
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
		Project:  c.Args().Get(0),
		Config:   c.String("file"),
		Settings: settings,
		Windows:  c.StringSlice("windows"),
		Attach:   c.Bool("attach"),
		Detach:   c.Bool("detach"),
		Debug:    c.Bool("debug"),
	}

}

func checkProjectArgument(c *cli.Context) {
	if c.Args().First() == "" {
		cli.ShowCommandHelp(c, c.Command.FullName())
		os.Exit(1)
	}
}

func getConfigPath(options Options) string {
	userConfigDir := filepath.Join(ExpandPath("~/"), ".config/smug")

	var configPath string
	if options.Config != "" {
		configPath = options.Config
	} else {
		configPath = filepath.Join(userConfigDir, options.Project+".yml")
	}
	return configPath
}

func getParsedConfig(options Options) Config {

	config, err := GetConfig(getConfigPath(options), options.Settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return config
}

func initSmug(options Options) Smug {

	configDir := filepath.Dir(getConfigPath(options))

	var logger *log.Logger

	if options.Debug {
		logFile, err := os.Create(filepath.Join(configDir, "smug.log"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		logger = log.New(logFile, "", 0)
	}

	commander := DefaultCommander{logger}
	tmux := Tmux{commander}
	smug := Smug{tmux, commander}

	return smug
}

func main() {

	context := CreateContext()

	// EXAMPLE: Override a template
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
EXAMPLES:
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

WEBSITE:
	https://github.com/ivaaaan/smug
{{if .Copyright }}COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
VERSION:
	{{.Version}}
{{end}}
`
	app := &cli.App{

		Usage:   fmt.Sprintf(`tmux session manager. Version %s`, version),
		Version: version,

		/* GLOBAL FLAGS */
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Usage:   "A custom path to a config file",
				Aliases: []string{"f"},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Print all commands to <CONFIG_DIR>/smug.log",
			},
		},

		Commands: []*cli.Command{

			/* START COMMAND */
			{
				Name:      "start",
				Usage:     "start project session",
				ArgsUsage: "[<project>] [<key>=<value>]...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "windows",
						Usage:   "List of windows to start. If session exists, those windows will be attached to current session.",
						Aliases: []string{"w"},
					},
					&cli.BoolFlag{
						Name:    "attach",
						Usage:   "Force switch client for a session",
						Aliases: []string{"a"},
					},

					&cli.BoolFlag{
						Name:    "detach",
						Usage:   "Detach tmux session. The same as -d flag in the tmux",
						Aliases: []string{"d"},
					},
				},
				Action: func(c *cli.Context) error {

					checkProjectArgument(c)

					options := initOptions(c)
					smug := initSmug(options)

					if len(options.Windows) == 0 {
						fmt.Println("Starting a new session...")
					} else {
						fmt.Println("Starting new windows...")
					}

					config := getParsedConfig(options)

					err := smug.Start(config, options, context)
					if err != nil {
						fmt.Println("Oops, an error occurred! Rolling back...")
						smug.Stop(config, options, context)
						os.Exit(1)
					}

					return nil
				},
			},

			/* STOP COMMAND */
			{
				Name:      "stop",
				Usage:     "stop project session",
				ArgsUsage: "[<project>] [<key>=<value>]...",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "windows",
						Usage:   "List of windows to start. If session exists, those windows will be attached to current session.",
						Aliases: []string{"w"},
					},
				},
				Action: func(c *cli.Context) error {

					checkProjectArgument(c)

					options := initOptions(c)
					smug := initSmug(options)

					if len(options.Windows) == 0 {
						fmt.Println("Terminating session...")
					} else {
						fmt.Println("Killing windows...")
					}

					config := getParsedConfig(options)
					err := smug.Stop(config, options, context)
					if err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
						os.Exit(1)
					}

					return nil
				},
			},

			/* LIST COMMAND */
			{
				Name:  "list",
				Usage: "list available session configurations",
				Action: func(c *cli.Context) error {

					options := initOptions(c)
					configs, err := ListConfigs(filepath.Dir(getConfigPath(options)))

					if err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
						os.Exit(1)
					}

					fmt.Println(strings.Join(configs, "\n"))
					return nil
				},
			},

			/* PRINT COMMAND */
			{
				Name:  "print",
				Usage: "Print session configuration to stdout",
				Action: func(c *cli.Context) error {

					options := initOptions(c)
					smug := initSmug(options)

					config, err := smug.GetConfigFromSession(options, context)
					if err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
						os.Exit(1)
					}

					d, err := yaml.Marshal(&config)
					if err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
						os.Exit(1)
					}

					fmt.Println(string(d))
					return nil
				},
			},

			/* EDIT COMMAND */
			{
				Name:      "edit",
				Usage:     "edit project session configuration",
				Aliases:   []string{"new"},
				ArgsUsage: "[<project>] [<key>=<value>]...",
				Action: func(c *cli.Context) error {

					checkProjectArgument(c)

					options := initOptions(c)
					err := EditConfig(getConfigPath(options))

					if err != nil {
						fmt.Fprintf(os.Stderr, err.Error())
						os.Exit(1)
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
