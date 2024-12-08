package main

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type repoMsg struct {
	repo
}
type doneMsg struct{}
type errMsg error

func (r repo) Title() string {
	return r.path
}

func (r repo) Description() string {
	// NOTE: list renders single line only..
	return r.status
}

func (r repo) FilterValue() string {
	return r.path
}

type model struct {
	list     list.Model
	path     string
	sub      chan repo
	quitting bool
	err      error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		startRepoSearch(m.path, m.sub),
		waitForRepoStatus(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, openEditor(m.list.SelectedItem().(repo).path)
		case "u":
			return m, tea.Sequence(m.list.SetItems([]list.Item{}), m.Init())
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
	// Remove duplicate (status has it) + misaligned string coming from:
	// https://github.com/charmbracelet/bubbles/blob/178590b4469b2386726cff8da7c479615a746a94/list/list.go#L1220
	s := strings.Replace(m.list.View(), "No repositories.", "", 1)
	if m.quitting {
		s += "\n"
	}
	return s
}

func initialModel(wd string) tea.Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Dirty Repositories"
	l.SetStatusBarItemName("repository", "repositories")
	l.SetSpinner(spinner.MiniDot)
	l.ToggleSpinner()
	l.AdditionalShortHelpKeys = additionalShortHelpKeys
	return model{list: l, path: wd, sub: make(chan repo)}
}

func additionalShortHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update"),
		),
	}
}

func startRepoSearch(path string, sub chan repo) tea.Cmd {
	return func() tea.Msg {
		err := collectDirtyRepos(path, sub)
		if err != nil {
			return err
		}
		return doneMsg{}
	}
}

func waitForRepoStatus(sub chan repo) tea.Cmd {
	return func() tea.Msg {
		repo := <-sub
		return repoMsg(repoMsg{repo: repo})
	}
}

func openEditor(path string) tea.Cmd {
	c := exec.Command("code", path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return nil
	})
}
