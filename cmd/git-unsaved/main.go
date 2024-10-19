package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

type errMsg error

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n%s Scanning Git repositories...\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
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
			p := tea.NewProgram(initialModel())
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
