package cli

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jonnyarnold/fn-go/compiler"
	"github.com/jonnyarnold/fn-go/repl"
	"github.com/jonnyarnold/fn-go/vm"
)

type fnCli struct {
	app *cli.App
}

// CLI definition.
func buildCli() fnCli {
	app := cli.NewApp()
	app.Name = "fn"
	app.Usage = "A fun-ctional programming language!"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Runs the given filename as a a script.",
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				compiler.Run(fileName)
			},
		},

		{
			Name:    "repl",
			Aliases: []string{"i"},
			Usage:   "Starts an interactive REPL.",
			Action: func(c *cli.Context) {
				repl.Run()
			},
		},

		{
			Name:  "vm",
			Usage: "Directly run bytecode in the VM.",
			Action: func(c *cli.Context) {
				fmt.Println(vm.RunBytecode())
			},
		},
	}

	return fnCli{app: app}
}

// Passes control to the CLI.
func (fnCli fnCli) Run(args []string) {
	fnCli.app.Run(args)
}

// Use this outside!
var CLI = buildCli()
