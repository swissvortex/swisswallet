package model

import (
	. "vortex-wallet/enum"
)

type Arguments struct {
	Currency string `json:"currency"`
	Secret   string `json:"secret"`
	Salt     string `json:"salt"`
}

func (a *Arguments) GetCurrencyCode() int {
	return CurrencyCode[a.Currency]
}

func (a *Arguments) GetSecret() string {
	return a.Secret
}

func (a *Arguments) GetSalt() string {
	return a.Salt
}

func (a *Arguments) GetCurrencySecretByKdf(kdfType int) string {
	return a.Secret + string(rune(a.GetCurrencyCode()+kdfType))
}

func (a *Arguments) GetCurrencySaltByKdf(kdfType int) string {
	return a.Salt + string(rune(a.GetCurrencyCode()+kdfType))
}

func (a *Arguments) PassphraseIsEmpry() bool {
	if a.Secret == "" {
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
