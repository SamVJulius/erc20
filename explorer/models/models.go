package models

type Block struct {
	Number  uint64 `json:"number"`
	TxCount int    `json:"txcount"`
}

type Transaction struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
}

type TokenTransfer struct {
	TxHash string `json:"txhash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Value  string `json:"value"`
}
