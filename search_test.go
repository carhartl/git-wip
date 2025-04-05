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
	defer func() {
		_ = os.RemoveAll(d)
	}()
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

func TestExcludedDirectory(t *testing.T) {
	d, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		_ = os.RemoveAll(d)
	}()
	err = os.MkdirAll(d+"/.hidden/test", 0750)
	if err != nil {
		t.Error(err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = d + "/.hidden/test"
	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}
	err = os.MkdirAll(d+"/node_modules/test", 0750)
	if err != nil {
		t.Error(err)
	}
	cmd = exec.Command("git", "init")
	cmd.Dir = d + "/node_modules/test"
	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	output := make(chan repo)
	go func() {
		_ = collectDirtyRepos(d, output)
		defer close(output)
	}()

	detected := 0
	select {
	case result := <-output:
		// Closed channel returning the zero value of the underlying type, thus check non-emptyness..
		if result.path != "" {
			detected++
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}

	require.Equal(t, 0, detected)
}
