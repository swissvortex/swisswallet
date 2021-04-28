package constants

const LOGGING_LEVEL string = "error"
const AES_BLOCKSIZE_ERROR string = "AES BlockSize error: ciphertext too short"
const AES_PLAINTEXT_NOT_MULTIPLE_ERROR string = "Plaintext is not a multiple of the block size"

const GENERATE_MODE string = "generate"
const DECRYPT_MODE string = "decrypt"
const ENCRYPT_MODE string = "encrypt"

const RAW_OUTPUT string = "raw"
const MNEMONIC_OUTPUT string = "mnemonic"

const (
	UNKNOWN int = iota - 1
	TESTNET
	BITCOIN
	ETHEREUM
	LITECOIN
	MONERO
	COSMOS
	POLKADOT
)

var CurrencyCode = map[string]int{
	"unknown":  UNKNOWN,
	"testnet":  TESTNET,
	"bitcoin":  BITCOIN,
	"ethereum": ETHEREUM,
	"litecoin": LITECOIN,
	"monero":   MONERO,
	"cosmos":   COSMOS,
	"polkadot": POLKADOT,
}

const (
	ARGON2 int = iota
	SCRYPT
)
