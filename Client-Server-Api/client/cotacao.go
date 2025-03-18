package client

type Cotacao struct {
	UsdBrl string `json:"bid"`
}

func NewCotacao(dolarEmReais string) *Cotacao {
	return &Cotacao{
		UsdBrl: dolarEmReais,
	}
}
