package utils

import (
	fs "flag"
	"fmt"
	"os"
	"vortex-wallet/model"
)

type SimpleUtils interface {
	IsEmptyString(str string) bool
	IsEmptyArray(array []string) bool
	GetArguments() (*model.Arguments, []string)
	PrintHelp()
	ExitWithError(err error)
}

type simpleUtils struct{}

func NewSimpleUtils() SimpleUtils {
	return &simpleUtils{}
}

func (s *simpleUtils) IsEmptyString(str string) bool {
	if str == "" {
		return true
	} else {
		return false
	}
}

func (s *simpleUtils) IsEmptyArray(array []string) bool {
	if len(array) == 0 {
		return true
	} else {
		return false
	}
}

func (s *simpleUtils) GetArguments() (*model.Arguments, []string) {
	arguments := new(model.Arguments)

	fs.StringVar(&arguments.Currency, "currency", "ethereum", "Currency to use. Currently supported are testnet, bitcoin, ethereum, litecoin, monero, cosmos, polkadot")
	fs.StringVar(&arguments.Secret, "secret", "", "Vortex-wallet secret")
	fs.StringVar(&arguments.Salt, "salt", "", "Vortex-wallet salt")
	fs.Parse()

	return arguments, fs.Args()
}

func (s *simpleUtils) PrintHelp() {
	fmt.Printf("Usage: %s [options...]\n", os.Args[0])
	fs.PrintDefaults()
	fmt.Printf("Required: secret, salt\n")
	os.Exit(1)
}

func (s *simpleUtils) ExitWithError(err error) {
	fmt.Printf("%s\n", err)
	os.Exit(1)
}
