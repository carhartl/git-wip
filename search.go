package main

import (
	"bytes"
	"io/fs"
	"os/exec"
	"path/filepath"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
)

type repoMsg struct {
	repo repo
}
type repo struct {
	path   string
	status string
}

var excludeDirs = regexp.MustCompile(`.+/(\..+|node_modules)`) // Skip hidden directories (incl. .git) and node_modules

func getRepos(path string, sub chan repoMsg) tea.Cmd {
	return func() tea.Msg {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && excludeDirs.MatchString(path) {
				if filepath.Base(path) == ".git" {
					repopath := filepath.Dir(path)

					var buf = new(bytes.Buffer)
					cmd := exec.Command("git", "status", "--porcelain=v2", "--show-stash")
					cmd.Stdout = buf
					cmd.Dir = repopath
					err = cmd.Run()
					if err != nil {
						return err
					}

					gi := GitInfo{}
					gi.Parse(buf)

					if !gi.IsClean() {
						sub <- repoMsg{repo: repo{path: repopath, status: gi.Summary()}}
					}
				}
				return fs.SkipDir
			}
			return nil
		})
		if err != nil {
			return err
		}
		return doneMsg{}
	}
}

func waitForRepoStatus(sub chan repoMsg) tea.Cmd {
	return func() tea.Msg {
		return repoMsg(<-sub)
	}
}
