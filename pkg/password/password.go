// Package password provides a library for generating random password strings.
//
//	generator := password.NewPassword()
//	generator.AllowRepeat = true
//	pass, err := generator.Generate(16)
//	if err != nil {
//	   log.Fatal(err)
//	}
//	log.Printf(pass)
//
// Author: Lugh
package password

import (
	"crypto/rand"
	"math/big"
	"slices"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLettersDefault = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLettersDefault = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	DigitsDefault = "0123456789"

	// Symbols is the list of symbols.
	SymbolsDefault = "!@#$%^&*()_+"

	// Don't allow repeat by default
	AllowRepeatDefault = false
)

// Represent the configuration for the generator.
type Password struct {
	LowerLetters string
	UpperLetters string
	Digits       string
	Symbols      string
	AllowRepeat  bool
}

// Creates a new password generator with default configuration.
// To specified a configuration modify directly the exported field.
//
// Default values:
//   - "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+"
//   - Don't allow repeat by default.
func NewPassword() *Password {
	return &Password{
		LowerLetters: LowerLettersDefault,
		UpperLetters: UpperLettersDefault,
		Digits:       DigitsDefault,
		Symbols:      SymbolsDefault,
		AllowRepeat:  AllowRepeatDefault,
	}
}

// Generates a password with the given requirements. length is the
// total number of characters in the password.
func (p *Password) Generate(length int) (string, error) {

	var availableChars string = p.UpperLetters + p.LowerLetters + p.Digits + p.Symbols

	chars := []rune(availableChars)

	passwd := make([]rune, length)

	for i := 0; i < length; i++ {
	NoRepeat:
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		if !p.AllowRepeat && slices.Contains(passwd, chars[index.Int64()]) {
			goto NoRepeat
		}
		passwd[i] = chars[index.Int64()]
	}

	return string(passwd), nil
}
