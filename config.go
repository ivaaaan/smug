package main

import "gopkg.in/yaml.v2"

type Pane struct {
	Root     string   `yaml:"root"`
	Type     string   `yaml:"type"`
	Commands []string `yaml:"commands"`
}

type Window struct {
	Name        string   `yaml:"name"`
	Root        string   `yaml:"root"`
	BeforeStart []string `yaml:"before_start"`
	Panes       []Pane   `yaml:"panes"`
	Commands    []string `yaml:"commands"`
	Manual      bool
}

type Config struct {
	Session     string   `yaml:"session"`
	Root        string   `yaml:"root"`
	BeforeStart []string `yaml:"before_start"`
	Stop        []string `yaml:"stop"`
	Windows     []Window `yaml:"windows"`
}

func ParseConfig(data string) (*Config, error) {
	c := Config{}

	err := yaml.Unmarshal([]byte(data), &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
