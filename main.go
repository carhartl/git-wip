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
  -h, --help`

	app := &cli.App{
		Name:  "git-wip",
		Usage: "Find all your dirty Git repositories",
		Action: func(ctx *cli.Context) error {
			dir := ctx.Args().First()
			if dir == "" {
				dir, _ = os.Getwd() // Default: current directory
			}
			p := tea.NewProgram(initialModel(dir), tea.WithAltScreen())
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
