package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/urfave/cli/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var excludeDirs = regexp.MustCompile(`.+/(\..+|node_modules)`) // Skip hidden directories (incl. .git) and node_modules

type repoMsg struct {
	repo repo
}
type doneMsg struct{}
type errMsg error

type repo struct {
	path string
}

func (r repo) Title() string {
	return r.path
}

func (r repo) Description() string {
	return "Foo" // TODO: Gather Git status here, separated by \n
}

func (r repo) FilterValue() string {
	return r.path
}

type model struct {
	list     list.Model
	path     string
	sub      chan repoMsg
	quitting bool
	err      error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		getRepos(m.path, m.sub),
		waitForRepoStatus(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// `to do script "..." in window 1` won't open a new Terminal window but it won't cd,
			// I assume because git-unsaved is still running...
			c := exec.Command("osascript", "-e", "tell application \"Terminal\" to do script \"cd "+m.list.SelectedItem().(repo).path+"\"", ">/dev/null")
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return nil
			})
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case repoMsg:
		return m, tea.Sequence(
			m.list.InsertItem(len(m.list.Items()), msg.repo),
			waitForRepoStatus(m.sub), // wait for next found repo
		)
	case doneMsg:
		m.list.StopSpinner()
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	s := m.list.View()
	if m.quitting {
		s += "\n"
	}
	return s
}

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

					status, err := w.Status()
					if err != nil {
						return err
					}

					if !status.IsClean() {
						sub <- repoMsg{repo: repo{path: repopath}}
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

func main() {
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

Usage:
  git-unsaved [path] [flags]

Examples:
  git-unsaved
  git-unsaved /path/to/directory
  git unsaved

Flags:
  -h, --help
  -v, --version`

	app := &cli.App{
		Name:    "git-unsaved",
		Usage:   "Find all your dirty Git repositories",
		Version: "v0.0.1",
		Action: func(ctx *cli.Context) error {
			l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
			l.Title = "Dirty Repositories"
			l.SetSpinner(spinner.MiniDot)
			l.ToggleSpinner()
			m := model{list: l, path: ctx.Args().First(), sub: make(chan repoMsg)}
			p := tea.NewProgram(m, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
