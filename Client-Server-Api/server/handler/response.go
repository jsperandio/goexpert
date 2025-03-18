package handler

import "errors"

var (
	ErrDatabaseTimeout              = errors.New("database context timeout")
	ErrExternalRequestTimeout       = errors.New("external request context timeout")
	ErrExternalRequestReturnedError = errors.New("external request returned error")
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type CotacaoResponse struct {
	UsdBrl UsdBrl `json:"USDBRL"`
}

type UsdBrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}
