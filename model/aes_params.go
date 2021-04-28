package model

type AESParams struct {
	EncryptionKey []byte `json:"key"`
	Input         []byte `json:"input"`
}

func (a *AESParams) GetIV() []byte {
	var iv []byte

	for i := 0; i < len(a.EncryptionKey)/2; i++ {
		iv = append(iv, a.EncryptionKey[i]^a.Input[i]^a.EncryptionKey[16+i]^a.Input[16+i])
	}

	return iv
}
