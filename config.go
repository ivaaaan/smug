package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Pane struct {
	Root     string   `yaml:"root,omitempty"`
	Type     string   `yaml:"type,omitempty"`
	Commands []string `yaml:"commands"`
}

type Window struct {
	Name        string   `yaml:"name"`
	Root        string   `yaml:"root,omitempty"`
	BeforeStart []string `yaml:"before_start"`
	Panes       []Pane   `yaml:"panes"`
	Commands    []string `yaml:"commands"`
	Layout      string   `yaml:"layout"`
	Manual      bool     `yaml:"manual,omitempty"`
}

type Config struct {
	Session                   string            `yaml:"session"`
	Env                       map[string]string `yaml:"env"`
	Root                      string            `yaml:"root"`
	BeforeStart               []string          `yaml:"before_start"`
	Stop                      []string          `yaml:"stop"`
	Windows                   []Window          `yaml:"windows"`
	RebalanceWindowsThreshold int               `yaml:"rebalance_panes_after"`
}

func EditConfig(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func GetConfigPathToProject(project string) string {
	return filepath.Join(GetUserConfigDir(), project+".yml")
}

// TODO make this return a slice of all config files
func GetConfigPathToAllProjects(project string) string {
	return filepath.Join(GetUserConfigDir(), project+".yml")
}

func GetUserConfigDir() string {
	return filepath.Join(ExpandPath("~/"), ".config/smug")
}

func GetConfig(path string, tmux Tmux, settings map[string]string, project string) (Config, error) {
	project = strings.TrimSpace(project)

	// If project is not defined in the arg, see if we can look it up from a current tmux session
	if len(project) == 0 {
		fmt.Println("NO PROJECT DEFINED")
		sessionName, _ := tmux.SessionName()
		fmt.Println("using " + sessionName)

		path = GetConfigPathToProject(sessionName)
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := string(f)

	return ParseConfig(config, settings)
}

func ParseConfig(data string, settings map[string]string) (Config, error) {
	data = os.Expand(data, func(v string) string {
		if val, ok := settings[v]; ok {
			return val
		}

		if val, ok := os.LookupEnv(v); ok {
			return val
		}

		return v
	})

	c := Config{}

	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func ListConfigs(dir string) ([]string, error) {
	var result []string
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return result, err
	}

	for _, file := range files {
		fileExt := path.Ext(file.Name())
		if fileExt != ".yml" {
			continue
		}
		result = append(result, strings.TrimSuffix(file.Name(), fileExt))
	}

	return result, nil
}
