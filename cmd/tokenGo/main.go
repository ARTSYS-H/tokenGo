package main

import (
	"fmt"
	"os"

	"github.com/ARTSYS-H/tokenGo/internal/passwordcli"
)

const help = `tokenGo is a collection of command to generates token.

Usage:

	tokenGo <command> [arguments]

The commands are:

	password	generate a password string

Use "tokenGo help <command>" for more information about a command.`

type SubCommand interface {
	Init([]string) error
	Run() error
	GetName() string
}

func root(args []string) error {
	if len(args) <= 1 {
		fmt.Println(help)
		return nil
	}

	commands := []SubCommand{
		passwordcli.NewPasswordCommand(),
	}

	subcommand := args[1]

	if subcommand == "help" {
		fmt.Println(help)
		return nil
	}

	for _, cmd := range commands {
		if cmd.GetName() == subcommand {
			err := cmd.Init(args[2:])
			if err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("%s: Unknown command\nRun 'help' for usage.", subcommand)
}

func main() {
	err := root(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
