package repo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

type CryptoRepository interface {
	DecryptCBC(key, ciphertext []byte) ([]byte, error)
	EncryptCBC(key, plaintext []byte) ([]byte, error)
	Argon2Kdf(secret string, salt string) []byte
	ScryptKdf(secret string, salt string) ([]byte, error)
}

type cryptoRepository struct{}

func NewCryptoRepository() CryptoRepository {
	return &cryptoRepository{}
}

func (c *cryptoRepository) DecryptCBC(key []byte, ciphertext []byte) ([]byte, error) {
	var block cipher.Block
	var err error

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		fmt.Printf("ciphertext too short")
		return nil, err
	}

	//TODO: Set proper IV
	iv := bytes.Repeat([]byte("0"), 16)

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	plaintext := ciphertext

	return plaintext, nil
}

func (c *cryptoRepository) EncryptCBC(key []byte, plaintext []byte) ([]byte, error) {
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	fmt.Printf("CBC Key: %s\n", hex.EncodeToString(key))
	fmt.Printf("CBC IV: %s\n", hex.EncodeToString(iv))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, err
}

func (c *cryptoRepository) Argon2Kdf(secret string, salt string) []byte {
	return argon2.Key([]byte(secret), []byte(salt), 8, 256*1024, 4, 32)
}

func (c *cryptoRepository) ScryptKdf(secret string, salt string) ([]byte, error) {
	return scrypt.Key([]byte(secret), []byte(salt), 262144, 8, 1, 32)
}
