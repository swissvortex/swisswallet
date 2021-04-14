# Vortex-wallet

Vortex-wallet is a deterministic cryptocurrency wallet generator heavily based on [MindWallet](https://github.com/patcito/mindwallet) and [MemWallet](https://github.com/dvdbng/memwallet) but using argon2 and scrypt by default as hashing functions. It's similar to [WarpWallet](https://keybase.io/warp/), but it works for Testnet, Bitcoin, Ethereum, Litecoin, Monero, Cosmos and Polkadot. You never need to save or store your private key anywhere. Just pick a really good password - many random words, for example - and never use it for anything else.

Given the same Secret and Salt, Vortex-wallet will always generate the same address and private key, so you only need to remember your password to access your funds.

For more information on why this is safer than a regular brainwallet, see [WarpWallet](https://keybase.io/warp/)'s help, Vortex-wallet is a re-implementation of WarpWallet, but it works for other currencies thanks to MemWallet and MindWallet who make the initial code. WarpWallet and MemWallet use the same algorithm, so WarpWallet and MemWallet will generate the same Bitcoin address for a given Passphrase and salt.

This repo contains an implementation of Vortex-wallet in Golang.