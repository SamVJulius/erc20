package services

import (
	"context"

	"explorer/rpc"
	"explorer/models"
)

func GetLatestBlock() (*models.Block, error) {
	block, err := rpc.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &models.Block{
		Number:  block.Number().Uint64(),
		TxCount: len(block.Transactions()),
	}, nil
}