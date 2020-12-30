package main

import "os"

type Context struct {
	InsideTmuxSession bool
}

func CreateContext() *Context {
	_, tmux := os.LookupEnv("TMUX")
	insideTmuxSession := os.Getenv("TERM") == "screen" || tmux
	return &Context{insideTmuxSession}
}
