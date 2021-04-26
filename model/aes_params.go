package model

type AESParams struct {
	Key   []byte `json:"key"`
	Input []byte `json:"input"`
}

func (a *AESParams) GetIV() []byte {
	var iv []byte

	for i := 0; i < len(a.Key)/2; i++ {
		iv = append(iv, a.Key[i]^a.Input[i]^a.Key[16+i]^a.Input[16+i])
	}

	return iv
}
