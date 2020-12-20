package main

import "fmt"

type ShellError struct {
	Command string
	Err     error
}

func (e *ShellError) Error() string {
	return fmt.Sprintf("Cannot run %q. Error %v", e.Command, e.Err)
}
