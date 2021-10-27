// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Tx struct {
	Hash    string      `json:"hash"`
	Raw     string      `json:"raw"`
	Inputs  []*TxInput  `json:"inputs"`
	Outputs []*TxOutput `json:"outputs"`
}

type TxInput struct {
	Tx     *Tx       `json:"tx"`
	Index  int       `json:"index"`
	Output *TxOutput `json:"output"`
}

type TxOutput struct {
	Tx     *Tx    `json:"tx"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
	Script string `json:"script"`
}