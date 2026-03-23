package services

import (
	"context"

	"explorer/config"
	"explorer/models"
	"explorer/rpc"
)

func GetLatestBlock() (*models.Block, error) {
	block, err := rpc.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &models.Block{
		Number:  block.Number().Uint64(),
		Gas:     block.GasUsed(),
		Hash:    block.Hash().String(),
		TxCount: len(block.Transactions()),
	}, nil
}

func GetConfig() (*models.Config, error) {
	block, err := rpc.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	chainId, err := rpc.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &models.Config{
		RpcUrl:          config.RPC_URL,
		ChainId:         chainId.Uint64(),
		ContractAddress: config.CONTRACT_ADDRESS,
		LatestBlock:     block.Number().Uint64(),
	}, nil
}
