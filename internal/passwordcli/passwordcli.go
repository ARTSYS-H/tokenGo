package passwordcli

import (
	"flag"
	"fmt"

	"github.com/ARTSYS-H/tokenGo/pkg/password"
)

type PasswordCommand struct {
	fs             *flag.FlagSet
	generator      *password.Password
	passwordLength *int
}

func NewPasswordCommand() *PasswordCommand {
	gen := password.NewPassword()
	fs := flag.NewFlagSet("password", flag.ExitOnError)

	passSize := fs.Int("length", 16, "Choose the password length.")

	fs.BoolFunc("allowrepeat", "Allow repeat in the password. (default false)", func(s string) error {
		gen.AllowRepeat = true
		return nil
	})

	return &PasswordCommand{
		fs:             fs,
		generator:      gen,
		passwordLength: passSize,
	}
}

func (sb *PasswordCommand) GetName() string {
	return sb.fs.Name()
}

func (sb *PasswordCommand) Init(args []string) error {
	return sb.fs.Parse(args)
}

func (sb *PasswordCommand) Run() error {
	generatedPassword, err := sb.generator.Generate(*sb.passwordLength)
	fmt.Println(generatedPassword)
	return err
}
