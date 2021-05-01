package repo

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"swisswallet/logger"

	. "swisswallet/constants"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

type CryptoRepository interface {
	AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error)
	AesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error)
	Argon2Kdf(password string, salt string, difficulty string) ([]byte, error)
	ScryptKdf(password string, salt string, difficulty string) ([]byte, error)
	GetArgon2ParamsByDifficulty(difficulty string) (uint32, uint32, uint8, uint32, error)
	GetScryptParamsByDifficulty(difficulty string) (int, int, int, int, error)
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

func (c *cryptoRepository) Argon2Kdf(password string, salt string, difficulty string) ([]byte, error) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), password, salt, difficulty)

	time, memory, threads, keyLen, err := c.GetArgon2ParamsByDifficulty(difficulty)
	if err != nil {
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	argon2Key := argon2.IDKey([]byte(password), []byte(salt), time, memory, threads, keyLen)

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", argon2Key))
	return argon2Key, err
}

func (c *cryptoRepository) ScryptKdf(password string, salt string, difficulty string) ([]byte, error) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), password, salt, difficulty)

	N, r, p, keyLen, err := c.GetScryptParamsByDifficulty(difficulty)
	if err != nil {
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	scryptKey, err := scrypt.Key([]byte(password), []byte(salt), N, r, p, keyLen)
	if err != nil {
		c.logger.LogOnInternalErrorWithContext(c.logger.GetContext(), err)
		return nil, err
	}

	c.logger.LogOnExitWithContext(c.logger.GetContext(), fmt.Sprintf("%x", scryptKey), err)
	return scryptKey, err
}

func (c *cryptoRepository) GetArgon2ParamsByDifficulty(difficulty string) (uint32, uint32, uint8, uint32, error) {
	switch difficulty {
	case RIDICULOUSLY_STRONG_DIFFICULTY:
		return 128, 8192 * 1024, 4, 32, nil
	case SUPER_STRONG_DIFFICULTY:
		return 64, 4096 * 1024, 4, 32, nil
	case STRONG_DIFFICULTY:
		return 32, 2048 * 1024, 4, 32, nil
	case NORMAL_DIFFICULTY:
		return 16, 1024 * 1024, 4, 32, nil
	case LOW_DIFFICULTY:
		return 8, 512 * 1024, 4, 32, nil
	case MINIMUM_DIFFICULTY:
		return 4, 256 * 1024, 4, 32, nil
	default:
		err := errors.New(fmt.Sprintf(BAD_REQUEST_DIFFICULTY_ERROR))
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		return 32, 2048 * 1024, 4, 32, err
	}
}

func (c *cryptoRepository) GetScryptParamsByDifficulty(difficulty string) (int, int, int, int, error) {
	switch difficulty {
	case RIDICULOUSLY_STRONG_DIFFICULTY:
		return 8192 * 1024, 8, 1, 32, nil
	case SUPER_STRONG_DIFFICULTY:
		return 4096 * 1024, 8, 1, 32, nil
	case STRONG_DIFFICULTY:
		return 2048 * 1024, 8, 1, 32, nil
	case NORMAL_DIFFICULTY:
		return 1024 * 1024, 8, 1, 32, nil
	case LOW_DIFFICULTY:
		return 512 * 1024, 8, 1, 32, nil
	case MINIMUM_DIFFICULTY:
		return 256 * 1024, 8, 1, 32, nil
	default:
		err := errors.New(fmt.Sprintf(BAD_REQUEST_DIFFICULTY_ERROR))
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		return 2048 * 1024, 8, 1, 32, err
	}
}
