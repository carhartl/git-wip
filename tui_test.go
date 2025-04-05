package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func withUntrackedFile(d string) error {
	cmd := exec.Command("touch", "foo.txt")
	cmd.Dir = d
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func TestSearchOutput(t *testing.T) {
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
	err = withUntrackedFile(d)
	if err != nil {
		t.Error(err)
	}

	tm := teatest.NewTestModel(t, initialModel(d), teatest.WithInitialTermSize(300, 100))
	if err != nil {
		t.Error(err)
	}

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("1 file to commit"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestUpdate(t *testing.T) {
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

	tm := teatest.NewTestModel(t, initialModel(d), teatest.WithInitialTermSize(300, 100))
	if err != nil {
		t.Error(err)
	}

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("missing upstream"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	err = withUntrackedFile(d)
	if err != nil {
		t.Error(err)
	}

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("u"),
	})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("1 file to commit, missing upstream"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
