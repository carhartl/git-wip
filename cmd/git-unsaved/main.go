package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

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
	sub      chan repoMsg
	quitting bool
	err      error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		getRepos(m.sub),
		waitForRepoStatus(m.sub),
	)
}

// TODO: upon selection + enter, cd to directory!
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
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

func getRepos(sub chan repoMsg) tea.Cmd {
	return func() tea.Msg {
		err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// TODO: filter for git repos
				abspath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				sub <- repoMsg{repo: repo{path: abspath}}
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
		Action: func(*cli.Context) error {
			l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
			l.Title = "Repositories"
			l.SetSpinner(spinner.MiniDot)
			l.ToggleSpinner()
			m := model{list: l, sub: make(chan repoMsg)}
			// TODO: how to pass along path argument?
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
