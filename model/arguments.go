package model

import (
	. "vortex-wallet/constants"
)

type Arguments struct {
	Password   string `json:"password"`
	Salt       string `json:"salt"`
	Currency   string `json:"currency"`
	Difficulty string `json:"difficulty"`
	Mnemonic   string `json:"mnemonic"`
	Key        string `json:"key"`
	Language   string `json:"language"`
	Address    string `json:"address"`
	Output     string `json:"output"`
}

func (a *Arguments) GetCurrencyCode() int {
	return CurrencyCode[a.Currency]
}

func (a *Arguments) GetPassword() string {
	return a.Password
}

func (a *Arguments) GetSalt() string {
	return a.Salt
}

func (a *Arguments) GetDifficulty() string {
	return a.Difficulty
}

func (a *Arguments) GetMnemonic() string {
	return a.Mnemonic
}

func (a *Arguments) GetKey() string {
	return a.Key
}

func (a *Arguments) GetLanguage() string {
	return a.Language
}

func (a *Arguments) GetAddress() string {
	return a.Address
}

func (a *Arguments) GetOutput() string {
	return a.Output
}

func (a *Arguments) GetCurrencyPasswordByKdf(kdfType int) string {
	return a.Password + string(rune(a.GetCurrencyCode()+kdfType))
}

func (a *Arguments) GetCurrencySaltByKdf(kdfType int) string {
	return a.Salt + string(rune(a.GetCurrencyCode()+kdfType))
}

func (a *Arguments) PasswordIsEmpry() bool {
	if a.Password == "" {
		return true
	} else {
		return false
	}
}

func (a *Arguments) SaltIsEmpry() bool {
	if a.Salt == "" {
		return true
	} else {
		return false
	}
}

func (a *Arguments) MnemonicIsEmpty() bool {
	if a.Mnemonic == "" {
		return true
	} else {
		return false
	}
}

func (a *Arguments) KeyIsEmpty() bool {
	if a.Key == "" {
		return true
	} else {
		return false
	}
}

func (a *Arguments) AddressIsEmpty() bool {
	if a.Address == "" {
		return true
	} else {
		return false
	}
}
