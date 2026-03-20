package services

import (
	"context"
	"explorer/models"
	"explorer/rpc"
)

func GetLatestTransactions() ([]models.Transaction, error) {
	block, err := rpc.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var txs []models.Transaction

	for _, tx := range block.Transactions() {
		txs = append(txs, models.Transaction{
			Hash: tx.Hash().Hex(),
			Value: tx.Value().String(),
		})
	}

	return txs, nil
}