package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ARTSYS-H/tokenGo/internal/passwordcli"
)

//go:embed help.txt
var help string

type SubCommand interface {
	Init([]string) error
	Run() error
	GetName() string
	ShowHelp()
}

type Cli struct {
	name        string
	subcommands []SubCommand
	help        string
}

func NewCli(name, helpText string) *Cli {
	return &Cli{
		name: name,
		help: helpText,
	}
}

func (c *Cli) Add(cmds ...SubCommand) {
	c.subcommands = append(c.subcommands, cmds...)
}

func (c *Cli) Run(args []string) error {
	if len(args) <= 1 {
		fmt.Println(c.help)
		return nil
	}

	if args[1] == "help" {
		if len(args) <= 2 {
			fmt.Println(c.help)
			return nil
		}

		for _, cmd := range c.subcommands {
			if cmd.GetName() == args[2] {
				cmd.ShowHelp()
				return nil
			}
		}

		return fmt.Errorf("%s help %s: unknown topic help. Run '%s help'.", c.name, args[2], c.name)
	}

	for _, cmd := range c.subcommands {
		if cmd.GetName() == args[1] {
			err := cmd.Init(args[2:])
			if err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("%s %s: Unknown command\nRun '%s help'.", c.name, args[1], c.name)
}

func main() {
	command := NewCli("tokenGo", help)
	command.Add(passwordcli.NewPasswordCommand())
	err := command.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
