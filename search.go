package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

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

					r, err := git.PlainOpen(repopath)
					if err != nil {
						return err
					}

					w, err := r.Worktree()
					if err != nil {
						return err
					}

					// Required until https://github.com/go-git/go-git/issues/1210 is fixed
					addDefaultGitignoreToWorktree(w)

					status, err := w.Status() // => git status --porcelain
					if err != nil {
						return err
					}

					if !status.IsClean() {
						sub <- repoMsg{repo: repo{path: repopath, status: status.String()}}
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

func addDefaultGitignoreToWorktree(w *git.Worktree) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	//TODO: first try $XDG_CONFIG_HOME/git/ignore, then fall back to $HOME/.config/git/ignore
	f, err := os.Open(filepath.Join(home, ".config", "git", "ignore"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		ignorePattern := sc.Text()
		w.Excludes = append(w.Excludes, gitignore.ParsePattern(ignorePattern, nil))
	}
}
