package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	"vortex-wallet/model"
	repo "vortex-wallet/repository"
	"vortex-wallet/utils"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"github.com/vsergeev/btckeygenie/btckey"
)

type Service interface {
	GenerateWallet(arguments model.Arguments) error
	DecryptWallet(arguments model.Arguments) error
	EncryptWallet(arguments model.Arguments) error
	GenerateAESParams(arguments model.Arguments) (*model.AESParams, error)
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
		fmt.Printf("Ethereum Address: 0x%x\n", address)
		fmt.Printf("Mnemonic: %s\n", mnemonic)
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) DecryptWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	err := s.simpleUtils.CheckIfSupported(arguments.Output, s.simpleUtils.GetSupportedOutputs())
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
	if arguments.Address == "0x"+hex.EncodeToString(address) {
		fmt.Println("Private Key successfully decrypted")
	} else {
		fmt.Println("Private Key does not match the provided address")
	}

	if arguments.Output == RAW_OUTPUT {
		fmt.Printf("Decrypted Private Key: %x\n", decryptedKeyAsBytes)
	} else {
		mnemonic, err := bip39.NewMnemonic(decryptedKeyAsBytes)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
			return err
		}
		fmt.Printf("Decrypted Mnemonic: %s\n", mnemonic)
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) EncryptWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	err := s.simpleUtils.CheckIfSupported(arguments.Output, s.simpleUtils.GetSupportedOutputs())
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
	}

	var privateKey btckey.PrivateKey
	entropyAsBytes, err := hex.DecodeString(arguments.Key)
	privateKey.FromBytes(entropyAsBytes)
	publicKey := privateKey.ToBytesUncompressed()
	address := ethcrypto.Keccak256(publicKey[1:])[12:]
	arguments.Salt = "0x" + hex.EncodeToString(address)

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
	fmt.Printf("Ethereum Address: 0x%x\n", address)

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) GenerateAESParams(arguments model.Arguments) (*model.AESParams, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	params := new(model.AESParams)
	var err error

	params.EncryptionKey = s.cryptoRepository.Argon2Kdf(arguments.GetCurrencyPasswordByKdf(ARGON2), arguments.GetCurrencySaltByKdf(ARGON2), arguments.GetDifficulty())
	params.Input, err = s.cryptoRepository.ScryptKdf(arguments.GetCurrencyPasswordByKdf(SCRYPT), arguments.GetCurrencySaltByKdf(SCRYPT), arguments.GetDifficulty())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return nil, err
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), fmt.Sprintf("%x", params), err)
	return params, err
}
