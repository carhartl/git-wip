package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDetectGitRepo(t *testing.T) {
	d, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(d)
	cmd := exec.Command("git", "init")
	cmd.Dir = d
	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	output := make(chan repo)
	go func() {
		_ = collectDirtyRepos(d, output)
		defer close(output)
	}()

	select {
	case result := <-output:
		require.Equal(t, d, result.path)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}
