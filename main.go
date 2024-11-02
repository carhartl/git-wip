package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

Usage:
  git-wip [path] [flags]

Examples:
  git-wip
  git-wip /path/to/directory
  git wip

Flags:
  -h, --help
  -v, --version`

	app := &cli.App{
		Name:    "git-wip",
		Usage:   "Find all your dirty Git repositories",
		Version: "v0.0.1",
		Action: func(ctx *cli.Context) error {
			p := tea.NewProgram(initialModel(ctx.Args().First()), tea.WithAltScreen())
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
