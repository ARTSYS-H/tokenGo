package main

import (
	"fmt"
	"os"

	"github.com/ARTSYS-H/crow/pkg/crow"
	"github.com/ARTSYS-H/tokenGo/internal/passwordcli"
)

func main() {
	app := crow.New("tokenGo", "TokenGo is a collection of commands to generate token.")
	err := app.AddCommand(passwordcli.NewPasswordCommand(), "generate a passwod string")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = app.Execute(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
