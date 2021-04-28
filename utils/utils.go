package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	"vortex-wallet/model"
)

type SimpleUtils interface {
	GetArguments() (*model.Arguments, string, []string)
	PrintHelpModeAndExit()
	PrintHelpParamsAndExit(mode string)
	ExitWithError(err error)
	GetSupportedModes() []string
	GetSupportedOutputs() []string
	CheckIfSupported(str string, supportedStrArray []string) error
	IsEmptyString(str string) bool
	IsEmptyArray(array []string) bool
	StringInSlice(a string, list []string) bool
}

type simpleUtils struct {
	logger *logger.Logger
}

func NewSimpleUtils(logger *logger.Logger) SimpleUtils {
	return &simpleUtils{
		logger: logger,
	}
}

var supportedModes = []string{GENERATE_MODE, DECRYPT_MODE, ENCRYPT_MODE}
var supportedOutputs = []string{RAW_OUTPUT, MNEMONIC_OUTPUT}

var fs = flag.NewFlagSet("options", flag.ContinueOnError)

func (s *simpleUtils) GetArguments() (*model.Arguments, string, []string) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), nil)
	arguments := new(model.Arguments)

	if s.IsEmptyArray(os.Args[1:]) {
		err := errors.New("Missing operation mode")
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		s.PrintHelpModeAndExit()
	}

	mode := os.Args[1]
	err := s.CheckIfSupported(mode, supportedModes)
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		s.ExitWithError(err)
	}

	fs.StringVar(&arguments.Password, "p", "", "Vortex-wallet password")
	fs.StringVar(&arguments.Salt, "s", "", "Vortex-wallet salt")
	fs.StringVar(&arguments.Currency, "c", "ethereum", "Currency to use. Currently supported are [testnet|bitcoin|ethereum|litecoin|monero|cosmos|polkadot")
	fs.StringVar(&arguments.Difficulty, "d", "strong", "Difficulty of the hashing algorithms. Currently supported are [minimum|low|normal|strong|super_strong|ridiculously_strong]")
	fs.StringVar(&arguments.Mnemonic, "m", "", "24 words mnemonic")
	fs.StringVar(&arguments.Key, "k", "", "Private key")
	fs.StringVar(&arguments.Language, "l", "english", "Mnemonic language [english|spanish|chinese_trad|chinese_simp|czech|french|italian|japanese|korean]")
	fs.StringVar(&arguments.Address, "a", "", "Currency address")
	fs.StringVar(&arguments.Output, "o", "mnemonic", "Output wallet format [mnemonic|raw]")
	fs.Parse(os.Args[2:])

	s.logger.LogOnExitWithContext(s.logger.GetContext(), arguments, mode, fs.Args())
	return arguments, mode, fs.Args()
}

func (s *simpleUtils) PrintHelpModeAndExit() {
	fmt.Println()
	fmt.Printf("Usage: %s %s [options...]\n", os.Args[0], supportedModes)
	PrintModes()
	os.Exit(1)
}

func (s *simpleUtils) PrintHelpParamsAndExit(mode string) {
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
	fmt.Println("- \"generate mnemonic\": vortex-wallet generate -p password -s salt")
	fmt.Println("- \"generate raw key\": vortex-wallet generate -o raw -p password -s salt")
	fmt.Println("- \"encrypt mnemonic\": vortex-wallet encrypt -m mnemonic -p password")
	fmt.Println("- \"encrypt raw key\": vortex-wallet encrypt -o raw -k privatekey -p password")
	fmt.Println("- \"decrypt mnemonic\": vortex-wallet decrypt -m mnemonic -p password -a address")
	fmt.Println("- \"decrypt raw key\": vortex-wallet decrypt -o raw -k privatekey -p password -a address")
	fmt.Println()
}

func (s *simpleUtils) GetSupportedModes() []string {
	return supportedModes
}

func (s *simpleUtils) GetSupportedOutputs() []string {
	return supportedOutputs
}

func (s *simpleUtils) CheckIfSupported(str string, supportedStrArray []string) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), str, supportedStrArray)

	var err error
	if !s.StringInSlice(str, supportedStrArray) {
		err = errors.New(fmt.Sprintf("Incorrect value: %s. Supported %s", str, supportedStrArray))
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

func (s *simpleUtils) StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
