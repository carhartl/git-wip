package main

import (
	"bytes"
	"io/fs"
	"os/exec"
	"path/filepath"
	"regexp"
)

type repo struct {
	path   string
	status string
}

var excludeDirs = regexp.MustCompile(`.+/(\..+|node_modules)`) // Skip hidden directories (incl. .git) and node_modules

func collectDirtyRepos(path string, sub chan repo) error {
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
				cmd := exec.Command("git", "status", "--porcelain=v2", "--show-stash", "--branch")
				cmd.Stdout = buf
				cmd.Dir = repopath
				err = cmd.Run()
				if err != nil {
					return err
				}

				gi := GitInfo{}
				gi.Parse(buf)

				if !gi.IsClean() {
					sub <- repo{path: repopath, status: gi.Summary()}
				}
			}
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
