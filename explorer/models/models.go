package models

type Block struct {
	Number  uint64 `json:"number"`
	Gas     uint64 `json:"gas"`
	Hash    string `json:"hash"`
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

type Config struct {
	RpcUrl          string `json:"rpcUrl"`
	ChainId         uint64 `json:"chainId"`
	ContractAddress string `json:"contractAddress"`
	LatestBlock     uint64 `json:"latestBlock"`
}
