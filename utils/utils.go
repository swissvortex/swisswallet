package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
	. "swisswallet/constants"
	"swisswallet/logger"
	"swisswallet/model"
)

type SimpleUtils interface {
	GetArguments() (*model.Arguments, string, []string)
	PrintHelpModeAndExit()
	PrintHelpParamsAndExit(mode string)
	ExitWithError(err error)
	GetSupportedModes() []string
	GetSupportedOutputs() []string
	GetSupportedLanguages() []string
	GetSupportedDifficulties() []string
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
var supportedLanguages = []string{ENGLISH_LANGUAGE, SPANISH_LANGUAGE, CHINESE_TRADITIONAL_LANGUAGE, CHINESE_SIMPLIFIED_LANGUAGE, CZECH_LANGUAGE, FRENCH_LANGUAGE, ITALIAN_LANGUAGE, JAPANESE_LANGUAGE, KOREAN_LANGUAGE}
var supportedDifficulties = []string{MINIMUM_DIFFICULTY, LOW_DIFFICULTY, NORMAL_DIFFICULTY, STRONG_DIFFICULTY, SUPER_STRONG_DIFFICULTY, RIDICULOUSLY_STRONG_DIFFICULTY}

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

	fs.StringVar(&arguments.Password, "p", "", "swisswallet password")
	fs.StringVar(&arguments.Salt, "s", "", "swisswallet salt")
	fs.StringVar(&arguments.Mnemonic, "m", "", "24 words mnemonic")
	fs.StringVar(&arguments.Key, "k", "", "Private key")
	fs.StringVar(&arguments.Address, "a", "", "Currency address")
	fs.StringVar(&arguments.Currency, "c", "ethereum", "Currency to use. Currently supported are [testnet|bitcoin|ethereum|litecoin|monero|cosmos|polkadot")
	fs.StringVar(&arguments.Difficulty, "d", SUPER_STRONG_DIFFICULTY, fmt.Sprintf("Difficulty of the hashing algorithms. Currently supported are %s", supportedDifficulties))
	fs.StringVar(&arguments.Language, "l", ENGLISH_LANGUAGE, fmt.Sprintf("Mnemonic language %s", supportedLanguages))
	fs.StringVar(&arguments.Output, "o", MNEMONIC_OUTPUT, fmt.Sprintf("Output wallet format %s", supportedOutputs))
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
	fmt.Println("- \"generate mnemonic\": swisswallet generate -p password -s salt")
	fmt.Println("- \"generate raw key\": swisswallet generate -o raw -p password -s salt")
	fmt.Println("- \"encrypt mnemonic\": swisswallet encrypt -m mnemonic -p password")
	fmt.Println("- \"encrypt raw key\": swisswallet encrypt -o raw -k privatekey -p password")
	fmt.Println("- \"decrypt mnemonic\": swisswallet decrypt -m mnemonic -p password -a address")
	fmt.Println("- \"decrypt raw key\": swisswallet decrypt -o raw -k privatekey -p password -a address")
	fmt.Println()
}

func (s *simpleUtils) GetSupportedModes() []string {
	return supportedModes
}

func (s *simpleUtils) GetSupportedOutputs() []string {
	return supportedOutputs
}

func (s *simpleUtils) GetSupportedLanguages() []string {
	return supportedLanguages
}

func (s *simpleUtils) GetSupportedDifficulties() []string {
	return supportedDifficulties
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
