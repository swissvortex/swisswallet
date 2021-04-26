package service

import (
	"fmt"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	"vortex-wallet/model"
	repo "vortex-wallet/repository"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	wordlist "github.com/tyler-smith/go-bip39/wordlists"
	"github.com/vsergeev/btckeygenie/btckey"
)

type Service interface {
	GenerateWallet(arguments model.Arguments) error
	GenerateAESParams(arguments model.Arguments) (*model.AESParams, error)
}

type service struct {
	cryptoRepository repo.CryptoRepository
	logger           *logger.Logger
}

func NewService(cryptoRepository repo.CryptoRepository, logger *logger.Logger) Service {
	return &service{
		cryptoRepository: cryptoRepository,
		logger:           logger,
	}
}

func (s *service) GenerateWallet(arguments model.Arguments) error {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	params, err := s.GenerateAESParams(arguments)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	entropy, err := s.cryptoRepository.AesDecrypt(params.Input, params.Key, params.GetIV())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}

	bip39.SetWordList(wordlist.English)
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return err
	}
	fmt.Printf("Generated mnemonic: %s\n", mnemonic)

	var privateKey btckey.PrivateKey
	privateKey.FromBytes(entropy)
	publicKey := privateKey.ToBytesUncompressed()
	address := ethcrypto.Keccak256(publicKey[1:])[12:]
	fmt.Printf("Ethereum Address: 0x%x\n", address)
	fmt.Printf("Private Key: %x\n", entropy)

	s.logger.LogOnExitWithContext(s.logger.GetContext(), err)
	return err
}

func (s *service) GenerateAESParams(arguments model.Arguments) (*model.AESParams, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), arguments)

	params := new(model.AESParams)
	var err error

	params.Key = s.cryptoRepository.Argon2Kdf(arguments.GetCurrencyPasswordByKdf(ARGON2), arguments.GetCurrencySaltByKdf(ARGON2), arguments.GetDifficulty())
	params.Input, err = s.cryptoRepository.ScryptKdf(arguments.GetCurrencyPasswordByKdf(SCRYPT), arguments.GetCurrencySaltByKdf(SCRYPT), arguments.GetDifficulty())
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return nil, err
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), fmt.Sprintf("%x", params), err)
	return params, err
}
