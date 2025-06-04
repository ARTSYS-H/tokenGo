package passwordcli

import (
	"fmt"

	"github.com/ARTSYS-H/tokenGo/pkg/password"
)

type Password struct {
	Length      int  `flag:"l" help:"Choose the password length."`
	AllowRepeat bool `help:"Allow repeat in the password"`
}

func NewPasswordCommand() *Password {
	return &Password{
		Length:      16,
		AllowRepeat: false,
	}
}

func (sb *Password) Run() error {
	generator := password.NewPassword()
	generator.AllowRepeat = sb.AllowRepeat
	password, err := generator.Generate(sb.Length)
	if err != nil {
		return err
	}
	fmt.Println(password)
	return nil
}
