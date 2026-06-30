package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Worktree struct {
	Path   string
	Branch string
}

// GitWorktrees lists the git worktrees of the repository containing root.
func GitWorktrees(commander Commander, root string) ([]Worktree, error) {
	cmd := exec.Command("git", "-C", root, "worktree", "list", "--porcelain")
	out, err := commander.Exec(cmd)
	if err != nil {
		return nil, err
	}

	return parseWorktrees(out), nil
}

// parseWorktrees parses the output of `git worktree list --porcelain`.
func parseWorktrees(out string) []Worktree {
	var worktrees []Worktree
	var current Worktree

	for line := range strings.SplitSeq(out, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			current = Worktree{Path: strings.TrimPrefix(line, "worktree ")}
		case strings.HasPrefix(line, "branch "):
			current.Branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		case line == "":
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
		}
	}

	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}

// FindWorktree returns the path of the worktree matching name, by either its
// branch name or its directory basename.
func FindWorktree(worktrees []Worktree, name string) (string, error) {
	for _, wt := range worktrees {
		if wt.Branch == name || filepath.Base(wt.Path) == name {
			return wt.Path, nil
		}
	}

	return "", fmt.Errorf("worktree %q not found", name)
}
