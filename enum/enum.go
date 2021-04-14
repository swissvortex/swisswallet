package enum

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
