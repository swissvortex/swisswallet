package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	. "swisswallet/constants"
	"swisswallet/logger"
	"swisswallet/model"
	repo "swisswallet/repository"
	"swisswallet/utils"

	"github.com/ethereum/go-ethereum/accounts"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	wordlist "github.com/tyler-smith/go-bip39/wordlists"
	"github.com/vsergeev/btckeygenie/btckey"
)

type Service interface {
	GenerateWallet(arguments model.Arguments) error
	DecryptWallet(arguments model.Arguments) error
	EncryptWallet(arguments model.Arguments) error
	GenerateAESParams(arguments model.Arguments) (*model.AESParams, error)
	ChangeMnemonicLanguageIfSupported(language string) error
}

type service struct {
	cryptoRepository repo.CryptoRepository
	simpleUtils      utils.SimpleUtils
	logger           *logger.Logger
}

func NewService(cryptoRepository repo.CryptoRepository, simpleUtils utils.SimpleUtils, logger *logger.Logger) Service {
	return &service{
		cryptoRepository: cryptoRepository,
		simpleUtils:      simpleUtils,
		logger:           logger,
	}
}

func (s *service) GenerateWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	err := s.simpleUtils.CheckIfSupported(arguments.Output, s.simpleUtils.GetSupportedOutputs())
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	err = s.ChangeMnemonicLanguageIfSupported(arguments.Language)
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	params, err := s.GenerateAESParams(arguments)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	entropy, err := s.cryptoRepository.AesDecrypt(params.Input, params.EncryptionKey, params.GetIV())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	var privateKey btckey.PrivateKey
	privateKey.FromBytes(entropy)
	publicKey := privateKey.ToBytesUncompressed()
	address := ethcrypto.Keccak256(publicKey[1:])[12:]

	if arguments.Output == RAW_OUTPUT {
		fmt.Printf("Ethereum Address: 0x%x\n", address)
		fmt.Printf("Private Key: %x\n", entropy)
	} else {
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		wallet, err := hdwallet.NewFromMnemonic(mnemonic)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
		account, err := wallet.Derive(path, false)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		fmt.Printf("Ethereum Address: %s\n", account.Address.Hex())
		fmt.Printf("Mnemonic: %s\n", mnemonic)
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) DecryptWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)
	var mnemonic string

	err := s.simpleUtils.CheckIfSupported(arguments.Output, s.simpleUtils.GetSupportedOutputs())
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	err = s.ChangeMnemonicLanguageIfSupported(arguments.Language)
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	if (arguments.MnemonicIsEmpty() && arguments.KeyIsEmpty()) || arguments.AddressIsEmpty() {
		err := errors.New("Private Key or Mnemonic, and address are required in decryption mode")
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	} else if arguments.KeyIsEmpty() && !arguments.MnemonicIsEmpty() {
		entropy, err := bip39.EntropyFromMnemonic(arguments.Mnemonic)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		arguments.Key = hex.EncodeToString(entropy)
	}

	arguments.Salt = arguments.Address
	params, err := s.GenerateAESParams(arguments)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	encryptedKeyAsBytes, err := hex.DecodeString(arguments.Key)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	decryptedKeyAsBytes, err := s.cryptoRepository.AesDecrypt(encryptedKeyAsBytes, params.EncryptionKey, params.GetIV())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	var privateKey btckey.PrivateKey
	privateKey.FromBytes(decryptedKeyAsBytes)
	publicKey := privateKey.ToBytesUncompressed()
	address := ethcrypto.Keccak256(publicKey[1:])[12:]
	mnemonic, err = bip39.NewMnemonic(decryptedKeyAsBytes)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	if arguments.Output == MNEMONIC_OUTPUT {
		wallet, err := hdwallet.NewFromMnemonic(mnemonic)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
		account, err := wallet.Derive(path, false)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		if arguments.Address == account.Address.Hex() {
			fmt.Println("Private Key successfully decrypted")
			fmt.Printf("Decrypted Mnemonic: %s\n", mnemonic)
		} else {
			fmt.Println("Private Key does not match the provided address, or wrong output mode")
		}
	} else if arguments.Output == RAW_OUTPUT {
		if arguments.Address == "0x"+hex.EncodeToString(address) {
			fmt.Println("Private Key successfully decrypted")
			fmt.Printf("Decrypted Private Key: %x\n", decryptedKeyAsBytes)
		} else {
			fmt.Println("Private Key does not match the provided address, or wrong output mode")
		}
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) EncryptWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)
	var account accounts.Account
	var privateKey btckey.PrivateKey
	var entropyAsBytes, address []byte

	err := s.simpleUtils.CheckIfSupported(arguments.Output, s.simpleUtils.GetSupportedOutputs())
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	err = s.ChangeMnemonicLanguageIfSupported(arguments.Language)
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	if arguments.MnemonicIsEmpty() && arguments.KeyIsEmpty() {
		err := errors.New("Private Key or Mnemonic are required in encryption mode")
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	} else if arguments.KeyIsEmpty() && !arguments.MnemonicIsEmpty() {
		entropy, err := bip39.EntropyFromMnemonic(arguments.Mnemonic)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		arguments.Key = hex.EncodeToString(entropy)
		wallet, err := hdwallet.NewFromMnemonic(arguments.Mnemonic)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
		account, err = wallet.Derive(path, false)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		entropyAsBytes, err = hex.DecodeString(arguments.Key)
		arguments.Salt = account.Address.Hex()
	} else {
		entropyAsBytes, err = hex.DecodeString(arguments.Key)
		privateKey.FromBytes(entropyAsBytes)
		publicKey := privateKey.ToBytesUncompressed()
		address = ethcrypto.Keccak256(publicKey[1:])[12:]
		arguments.Salt = "0x" + hex.EncodeToString(address)
	}

	params, err := s.GenerateAESParams(arguments)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	encryptedKeyAsBytes, err := s.cryptoRepository.AesEncrypt(entropyAsBytes, params.EncryptionKey, params.GetIV())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	if arguments.Output == RAW_OUTPUT {
		fmt.Printf("Encrypted Private Key: %x\n", encryptedKeyAsBytes)
	} else {
		mnemonic, err := bip39.NewMnemonic(encryptedKeyAsBytes)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		fmt.Printf("Encrypted Mnemonic: %s\n", mnemonic)
	}
	fmt.Printf("Ethereum Address: %s\n", arguments.Salt)

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) GenerateAESParams(arguments model.Arguments) (*model.AESParams, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	params := new(model.AESParams)
	var err error

	params.EncryptionKey, err = s.cryptoRepository.Argon2Kdf(arguments.GetCurrencyPasswordByKdf(ARGON2), arguments.GetCurrencySaltByKdf(ARGON2), arguments.GetDifficulty())
	if err != nil {
		s.logger.LogOnErrorWithContext(s.logger.GetContext(), err)
		return nil, err
	}

	params.Input, err = s.cryptoRepository.ScryptKdf(arguments.GetCurrencyPasswordByKdf(SCRYPT), arguments.GetCurrencySaltByKdf(SCRYPT), arguments.GetDifficulty())
	if err != nil {
		s.logger.LogOnErrorWithContext(s.logger.GetContext(), err)
		return nil, err
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), fmt.Sprintf("%x", params), err)
	return params, err
}

func (s *service) ChangeMnemonicLanguageIfSupported(language string) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), language)

	err := s.simpleUtils.CheckIfSupported(language, s.simpleUtils.GetSupportedLanguages())
	if err != nil {
		s.logger.LogOnBadRequestErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	if language != ENGLISH_LANGUAGE {
		switch language {
		case SPANISH_LANGUAGE:
			bip39.SetWordList(wordlist.Spanish)
		case CHINESE_TRADITIONAL_LANGUAGE:
			bip39.SetWordList(wordlist.ChineseTraditional)
		case CHINESE_SIMPLIFIED_LANGUAGE:
			bip39.SetWordList(wordlist.ChineseSimplified)
		case CZECH_LANGUAGE:
			bip39.SetWordList(wordlist.Czech)
		case FRENCH_LANGUAGE:
			bip39.SetWordList(wordlist.French)
		case ITALIAN_LANGUAGE:
			bip39.SetWordList(wordlist.Italian)
		case JAPANESE_LANGUAGE:
			bip39.SetWordList(wordlist.Japanese)
		case KOREAN_LANGUAGE:
			bip39.SetWordList(wordlist.Korean)
		}
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}
