package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"vortex-wallet/logger"
	"vortex-wallet/model"
)

type SimpleUtils interface {
	GetArguments() (*model.Arguments, string, []string)
	PrintHelpMode()
	PrintHelpParams(mode string)
	ExitWithError(err error)
	CheckMode(mode string) error
	IsEmptyString(str string) bool
	IsEmptyArray(array []string) bool
}

type simpleUtils struct {
	logger *logger.Logger
}

func NewSimpleUtils(logger *logger.Logger) SimpleUtils {
	return &simpleUtils{
		logger: logger,
	}
}

var fs = flag.NewFlagSet("options", flag.ContinueOnError)

func (s *simpleUtils) GetArguments() (*model.Arguments, string, []string) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), nil)
	arguments := new(model.Arguments)
	mode := os.Args[1]

	fs.StringVar(&arguments.Password, "p", "", "Vortex-wallet password")
	fs.StringVar(&arguments.Salt, "s", "", "Vortex-wallet salt")
	fs.StringVar(&arguments.Currency, "c", "ethereum", "Currency to use. Currently supported are [testnet|bitcoin|ethereum|litecoin|monero|cosmos|polkadot")
	fs.StringVar(&arguments.Difficulty, "d", "normal", "Difficulty of the hashing algorithms. Currently supported are [normal|strong|super_strong|ridiculously_strong]")
	fs.StringVar(&arguments.Mnemonic, "m", "", "24 words mnemonic")
	fs.StringVar(&arguments.Key, "k", "", "Private key")
	fs.StringVar(&arguments.Language, "l", "english", "Mnemonic language [english|spanish|chinese_trad|chinese_simp|czech|french|italian|japanese|korean]")
	fs.StringVar(&arguments.Address, "a", "", "Currency address")
	fs.StringVar(&arguments.Output, "o", "mnemonic", "Output wallet format [mnemonic|key]")
	fs.Parse(os.Args[2:])

	s.logger.LogOnExitWithContext(s.logger.GetContext(), arguments, mode, fs.Args())
	return arguments, mode, fs.Args()
}

func (s *simpleUtils) PrintHelpMode() {
	fmt.Println()
	fmt.Printf("Usage: %s %s [options...]\n", os.Args[0], supportedModes)
	PrintModes()
	os.Exit(1)
}

func (s *simpleUtils) PrintHelpParams(mode string) {
	fmt.Println()
	fmt.Printf("Usage: %s %s [options...]\n", os.Args[0], mode)
	PrintModes()
	fs.PrintDefaults()
	os.Exit(1)
}

func (s *simpleUtils) ExitWithError(err error) {
	fmt.Printf("%s\n", err)
	os.Exit(1)
}

func PrintModes() {
	fmt.Println("Supported modes with required arguments:")
	fmt.Println("- \"generate mnemonic\": vortex-wallet genwallet -p password -s salt")
	fmt.Println("- \"generate private key\": vortex-wallet genwallet -o key -p password -s salt")
	fmt.Println("- \"encrypt mnemonic\": vortex-wallet encrypt -m mnemonic -p password")
	fmt.Println("- \"encrypt private key\": vortex-wallet encrypt -o key -k privatekey -p password")
	fmt.Println("- \"decrypt mnemonic\": vortex-wallet decrypt -m mnemonic -p password -a address")
	fmt.Println("- \"decrypt private key\": vortex-wallet decrypt -o key -k privatekey -p password -a address")
	fmt.Println()
}

const supportedModes string = "[genwallet|decrypt|encrypt]"

func (s *simpleUtils) CheckMode(mode string) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), mode)

	var err error
	switch mode {
	case "genwallet":
		break
	case "decrypt":
		break
	case "encrypt":
		break
	default:
		err = errors.New(fmt.Sprintf("Incorrect mode: %s. Supported modes %s", mode, supportedModes))
		s.logger.LogOnErrorWithContext(s.logger.GetContext(), err)
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
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
