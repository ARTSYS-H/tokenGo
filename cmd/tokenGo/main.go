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
		if len(args) <= 2 {
			fmt.Println(help)
			return nil
		}

		for _, cmd := range commands {
			if cmd.GetName() == args[2] {
				cmd.ShowHelp()
				return nil
			}
		}

		return fmt.Errorf("tokenGo help %s: unknown topic help. Run 'tokenGo help'.", args[2])
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

	return fmt.Errorf("tokenGo %s: Unknown command\nRun 'tokenGo help' for usage.", subcommand)
}

func main() {
	err := root(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
