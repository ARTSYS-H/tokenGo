package passwordcli

import (
	"flag"
	"fmt"

	"github.com/ARTSYS-H/tokenGo/pkg/password"
)

type PasswordCommand struct {
	Fs             *flag.FlagSet
	Generator      *password.Password
	PasswordLength *int
	Description    string
}

func NewPasswordCommand() *PasswordCommand {
	gen := password.NewPassword()
	fs := flag.NewFlagSet("password", flag.ExitOnError)

	passSize := fs.Int("l", 16, "Choose the password `length`.")

	fs.BoolFunc("allowrepeat", "Allow repeat in the password. (default false)", func(s string) error {
		gen.AllowRepeat = true
		return nil
	})

	return &PasswordCommand{
		Fs:             fs,
		Generator:      gen,
		PasswordLength: passSize,
		Description:    "generate a password string",
	}
}

func (sb *PasswordCommand) Run() error {
	generatedPassword, err := sb.Generator.Generate(*sb.PasswordLength)
	if err != nil {
		return err
	}
	fmt.Println(generatedPassword)
	return nil
}
