package constants

const LOGGING_LEVEL string = "error"

const AES_BLOCKSIZE_ERROR string = "AES BlockSize error: ciphertext too short"
const AES_PLAINTEXT_NOT_MULTIPLE_ERROR string = "Plaintext is not a multiple of the block size"

const BAD_REQUEST_DIFFICULTY_ERROR string = "Provided difficulty not supported"

const GENERATE_MODE string = "generate"
const DECRYPT_MODE string = "decrypt"
const ENCRYPT_MODE string = "encrypt"

const RAW_OUTPUT string = "raw"
const MNEMONIC_OUTPUT string = "mnemonic"

const ENGLISH_LANGUAGE string = "english"
const SPANISH_LANGUAGE string = "spanish"
const CHINESE_TRADITIONAL_LANGUAGE string = "chinese_trad"
const CHINESE_SIMPLIFIED_LANGUAGE string = "chinese_simp"
const CZECH_LANGUAGE string = "czech"
const FRENCH_LANGUAGE string = "french"
const ITALIAN_LANGUAGE string = "italian"
const JAPANESE_LANGUAGE string = "japanese"
const KOREAN_LANGUAGE string = "korean"

const MINIMUM_DIFFICULTY string = "minimum"
const LOW_DIFFICULTY string = "low"
const NORMAL_DIFFICULTY string = "normal"
const STRONG_DIFFICULTY string = "strong"
const SUPER_STRONG_DIFFICULTY string = "super_strong"
const RIDICULOUSLY_STRONG_DIFFICULTY string = "ridiculously_strong"

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
