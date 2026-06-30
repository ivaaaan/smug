package main

import "testing"

const worktreeListOutput = `worktree /home/ivan/dev/smug
HEAD 7ac59da0000000000000000000000000000000000
branch refs/heads/master

worktree /home/ivan/dev/smug-feature
HEAD abc1230000000000000000000000000000000000
branch refs/heads/feature-x

worktree /home/ivan/dev/smug-detached
HEAD def4560000000000000000000000000000000000
detached
`

func TestParseWorktrees(t *testing.T) {
	worktrees := parseWorktrees(worktreeListOutput)

	expected := []Worktree{
		{Path: "/home/ivan/dev/smug", Branch: "master"},
		{Path: "/home/ivan/dev/smug-feature", Branch: "feature-x"},
		{Path: "/home/ivan/dev/smug-detached", Branch: ""},
	}

	if len(worktrees) != len(expected) {
		t.Fatalf("expected %d worktrees, got %d", len(expected), len(worktrees))
	}

	for i, wt := range worktrees {
		if wt != expected[i] {
			t.Errorf("expected %+v, got %+v", expected[i], wt)
		}
	}
}

func TestFindWorktree(t *testing.T) {
	worktrees := parseWorktrees(worktreeListOutput)

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"feature-x", "/home/ivan/dev/smug-feature", false},    // by branch
		{"smug-feature", "/home/ivan/dev/smug-feature", false}, // by directory basename
		{"smug-detached", "/home/ivan/dev/smug-detached", false},
		{"missing", "", true},
	}

	for _, tt := range tests {
		got, err := FindWorktree(worktrees, tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("FindWorktree(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("FindWorktree(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
