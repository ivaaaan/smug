package main

import "os"

type Context struct {
	InsideTmuxSession bool
}

func CreateContext() *Context {
	insideTmuxSession := os.Getenv("TERM") == "screen"
	return &Context{insideTmuxSession}
}
