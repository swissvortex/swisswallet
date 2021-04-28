package repo

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"vortex-wallet/logger"

	. "vortex-wallet/constants"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

type CryptoRepository interface {
	AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error)
	AesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error)
	Argon2Kdf(password string, salt string, difficulty string) []byte
	ScryptKdf(password string, salt string, difficulty string) ([]byte, error)
	GetArgon2ParamsByDifficulty(difficulty string) (uint32, uint32, uint8, uint32)
	GetScryptParamsByDifficulty(difficulty string) (int, int, int, int)
}

type cryptoRepository struct {
	logger *logger.Logger
}

func NewCryptoRepository(logger *logger.Logger) CryptoRepository {
	return &cryptoRepository{
		logger: logger,
	}
}

func (c *cryptoRepository) AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), fmt.Sprintf("%x, %x, %x", ciphertext, key, iv))

	block, err := aes.NewCipher(key)
	if err != nil {
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		err = errors.New(AES_BLOCKSIZE_ERROR)
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plaintext, ciphertext)

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", plaintext), err)
	return plaintext, err
}

func (c *cryptoRepository) AesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), fmt.Sprintf("%x, %x, %x", plaintext, key, iv))

	if len(plaintext)%aes.BlockSize != 0 {
		err := errors.New(AES_PLAINTEXT_NOT_MULTIPLE_ERROR)
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", ciphertext), err)
	return ciphertext, err
}

func (c *cryptoRepository) Argon2Kdf(password string, salt string, difficulty string) []byte {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), password, salt, difficulty)

	time, memory, threads, keyLen := c.GetArgon2ParamsByDifficulty(difficulty)
	argon2Key := argon2.IDKey([]byte(password), []byte(salt), time, memory, threads, keyLen)

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", argon2Key))
	return argon2Key
}

func (c *cryptoRepository) ScryptKdf(password string, salt string, difficulty string) ([]byte, error) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), password, salt, difficulty)

	N, r, p, keyLen := c.GetScryptParamsByDifficulty(difficulty)
	scryptKey, err := scrypt.Key([]byte(password), []byte(salt), N, r, p, keyLen)
	if err != nil {
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", scryptKey), err)
	return scryptKey, err
}

func (c *cryptoRepository) GetArgon2ParamsByDifficulty(difficulty string) (uint32, uint32, uint8, uint32) {
	switch difficulty {
	case "ridiculously_strong":
		return 128, 8192 * 1024, 4, 32
	case "super_strong":
		return 64, 4096 * 1024, 4, 32
	case "strong":
		return 32, 2048 * 1024, 4, 32
	case "normal":
		return 16, 1024 * 1024, 4, 32
	case "low":
		return 8, 512 * 1024, 4, 32
	case "minimum":
		return 4, 256 * 1024, 4, 32
	default:
		return 32, 2048 * 1024, 4, 32
	}
}

func (c *cryptoRepository) GetScryptParamsByDifficulty(difficulty string) (int, int, int, int) {
	switch difficulty {
	case "ridiculously_strong":
		return 8192 * 1024, 8, 1, 32
	case "super_strong":
		return 4096 * 1024, 8, 1, 32
	case "strong":
		return 2048 * 1024, 8, 1, 32
	case "normal":
		return 1024 * 1024, 8, 1, 32
	case "low":
		return 512 * 1024, 8, 1, 32
	case "minimum":
		return 256 * 1024, 8, 1, 32
	default:
		return 2048 * 1024, 8, 1, 32
	}
}
