package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/ARTSYS-H/tokenGo/internal/passwordcli"
)

type SubCommand interface {
	Run() error
}

type Cli struct {
	name        string
	subcommands []SubCommand
	args        []string
	description string
}

func (c *Cli) AddCommand(cmds ...SubCommand) {
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

func inspectAndAccessDescriptionField(sub SubCommand) (string, error) {

	valueOfSub := reflect.ValueOf(sub)

	for valueOfSub.Kind() == reflect.Interface || valueOfSub.Kind() == reflect.Ptr {
		valueOfSub = valueOfSub.Elem()
	}

	if valueOfSub.Kind() == reflect.Struct {
		field := valueOfSub.FieldByName("Description")
		if field.IsValid() {
			return field.String(), nil
		} else {
			return "", nil
		}
	} else {
		return "", fmt.Errorf("Sub does not contain a valid struct")
	}
}

func (c *Cli) generateMainHelpMessage() (string, error) {
	var helpMessage string

	helpMessage += fmt.Sprintf("%s\n", c.description)
	helpMessage += "\n"
	helpMessage += "Usage:\n"
	helpMessage += "\n"
	helpMessage += fmt.Sprintf("\t%s <command> [arguments]\n", c.name)
	helpMessage += "\n"
	helpMessage += "The commands are:\n"
	helpMessage += "\n"

	for _, cmd := range c.subcommands {
		flagSetField, err := inspectAndAccessFlagSetField(cmd)
		if err != nil {
			return "", err
		}
		descriptionField, err := inspectAndAccessDescriptionField(cmd)
		if err != nil {
			return "", err
		}
		helpMessage += fmt.Sprintf("\t%-15s %s\n", flagSetField.Name(), descriptionField)
	}

	helpMessage += "\n"
	helpMessage += fmt.Sprintf("Use \"%s help <command>\" for more information about a command.\n", c.name)

	return helpMessage, nil
}

func (c *Cli) helpHandler() error {
	if len(c.args) <= 2 {
		helpMessage, err := c.generateMainHelpMessage()
		if err != nil {
			return err
		}
		fmt.Println(helpMessage)
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
			err := flagSetField.Parse(c.args[2:])
			if err != nil {
				return err
			}
			return cmd.Run()
		}
	}
	return fmt.Errorf("%s %s: Unknown command\nRun '%s help' for usage.", c.name, c.args[1], c.name)
}

func (c *Cli) Execute(args []string) error {

	c.args = append(c.args, args...)

	helpRegexp := regexp.MustCompile(`^-h|--help|-help|help$`)

	if len(args) <= 1 || helpRegexp.MatchString(args[1]) {
		return c.helpHandler()
	}

	return c.commandsHandler()
}

func main() {
	app := &Cli{
		name:        "tokenGo",
		description: "TokenGo is a collection of commands to generate token.",
	}
	passwordSubCommand := passwordcli.NewPasswordCommand()
	app.AddCommand(passwordSubCommand)
	err := app.Execute(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
