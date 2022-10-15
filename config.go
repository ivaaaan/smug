package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigNotFoundError struct {
	Project string
}

func (e ConfigNotFoundError) Error() string {
	return fmt.Sprintf("config not found for project %s", e.Project)
}

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
	SendKeysTimeout int               `yaml:"sendkeys_timeout"`
	Session         string            `yaml:"session"`
	Env             map[string]string `yaml:"env"`
	Root            string            `yaml:"root"`
	BeforeStart     []string          `yaml:"before_start"`
	Stop            []string          `yaml:"stop"`
	Windows         []Window          `yaml:"windows"`
}

func addDefaultEnvs(c *Config, path string) {
	c.Env["SMUG_SESSION"] = c.Session
	c.Env["SMUG_SESSION_CONFIG_PATH"] = path
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

func GetConfig(path string, settings map[string]string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := string(f)

	c, err := ParseConfig(config, settings)
	if err != nil {
		return nil, err
	}

	addDefaultEnvs(&c, path)

	return &c, err

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

	c := Config{
		Env: make(map[string]string),
	}

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
		if fileExt != ".yml" && fileExt != ".yaml" {
			continue
		}
		result = append(result, file.Name())
	}

	return result, nil
}

func FindConfig(dir, project string) (string, error) {
	configs, err := ListConfigs(dir)
	if err != nil {
		return "", err
	}

	for _, config := range configs {
		fileExt := path.Ext(config)
		if strings.TrimSuffix(config, fileExt) == project {
			return config, nil
		}
	}

	return "", ConfigNotFoundError{Project: project}
}
