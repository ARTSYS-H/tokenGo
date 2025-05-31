package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/ARTSYS-H/tokenGo/internal/passwordcli"
)

//go:embed help.txt
var help string

type SubCommand interface {
	Init([]string) error
	Run() error
}

type Cli struct {
	name        string
	subcommands []SubCommand
	args        []string
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

func inspectAndAccessFlagSetField(sub SubCommand) (*flag.FlagSet, error) {
	var flagSetField *flag.FlagSet

	valueOfSub := reflect.ValueOf(sub)

	for valueOfSub.Kind() == reflect.Interface || valueOfSub.Kind() == reflect.Ptr {
		valueOfSub = valueOfSub.Elem()
	}

	if valueOfSub.Kind() == reflect.Struct {

		for i := range valueOfSub.NumField() {
			field := valueOfSub.Field(i)

			if field.Type() == reflect.TypeOf(&flag.FlagSet{}) {
				flagSetField = field.Interface().(*flag.FlagSet)
			}
		}
	} else {
		return nil, fmt.Errorf("Sub does not contain a valid struct")
	}

	return flagSetField, nil
}

func (c *Cli) helpHandler() error {
	if len(c.args) <= 2 {
		fmt.Println(c.help)
		return nil
	}

	for _, cmd := range c.subcommands {
		flagSetField, err := inspectAndAccessFlagSetField(cmd)
		if err != nil {
			return err
		}
		if flagSetField.Name() == c.args[2] {
			flagSetField.Usage()
			return nil
		}
	}

	return fmt.Errorf("%s help %s: unknown help topic. Run '%s help'.", c.name, c.args[2], c.name)
}

func (c *Cli) commandsHandler() error {
	for _, cmd := range c.subcommands {
		flagSetField, err := inspectAndAccessFlagSetField(cmd)
		if err != nil {
			return err
		}
		if flagSetField.Name() == c.args[1] {
			err := cmd.Init(c.args[2:])
			if err != nil {
				return err
			}
			return cmd.Run()
		}
	}
	return fmt.Errorf("%s %s: Unknown command\nRun '%s help'.", c.name, c.args[1], c.name)
}

func (c *Cli) Run(args []string) error {

	c.args = append(c.args, args...)

	if len(args) <= 1 || args[1] == "help" {
		return c.helpHandler()
	}

	return c.commandsHandler()
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
