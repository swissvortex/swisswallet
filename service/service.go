package service

import (
	"fmt"
	. "vortex-wallet/enum"
	"vortex-wallet/model"
	repo "vortex-wallet/repository"
)

type Service interface {
	GenerateWallet(arguments model.Arguments) error
}

type service struct {
	cryptoRepository repo.CryptoRepository
}

func NewService(cryptoRepository repo.CryptoRepository) Service {
	return &service{
		cryptoRepository: cryptoRepository,
	}
}

func (s *service) GenerateWallet(arguments model.Arguments) error {
	argon2Key := s.cryptoRepository.Argon2Kdf(arguments.GetCurrencySecretByKdf(ARGON2), arguments.GetCurrencySaltByKdf(ARGON2))
	scryptKey, err := s.cryptoRepository.ScryptKdf(arguments.GetCurrencySecretByKdf(SCRYPT), arguments.GetCurrencySaltByKdf(SCRYPT))

	if err != nil {
		panic(err)
	}

	result, _ := s.cryptoRepository.DecryptCBC(argon2Key, scryptKey)
	fmt.Printf("Wallet seed: %x\n", result)

	return nil
}
