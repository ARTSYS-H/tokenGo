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
	"fmt"
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

	// Define the possible characters
	var availableChars string = p.UpperLetters + p.LowerLetters + p.Digits + p.Symbols

	chars := []rune(availableChars)

	// no infinity loop if length exceeds available runes
	if len(chars) < length && !p.AllowRepeat {
		return "", fmt.Errorf("length exceeds available runes and repeats are not allowed")
	}

	// Create a rune slice for the password
	passwd := make([]rune, length)

	// Fill the slice with random characters
	for i := 0; i < length; i++ {
	NoRepeat:
		// Generate a random index
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		// if a repetition is found retry to generate a random index
		if !p.AllowRepeat && slices.Contains(passwd, chars[index.Int64()]) {
			goto NoRepeat
		}
		// Add the corresponding character to the random index
		passwd[i] = chars[index.Int64()]
	}

	// Return the password as a string
	return string(passwd), nil
}
