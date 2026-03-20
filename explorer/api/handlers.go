package api

import (
	"encoding/json"
	"explorer/services"
	"net/http"
)


func BlockHandler(w http.ResponseWriter, r *http.Request) {
	block, err := services.GetLatestBlock()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 
	}

	json.NewEncoder(w).Encode(block)
}

func TxHandler(w http.ResponseWriter, r *http.Request) {
	txs, err := services.GetLatestTransactions() 
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 
	}

	json.NewEncoder(w).Encode(txs)
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	transfers, err := services.GetTokenTransfers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(transfers)
}