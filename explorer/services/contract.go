package services

import (
	"context"
	"math/big"

	"explorer/config"
	"explorer/models"
	"explorer/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetContractTransfers() ([]models.TokenTransfer, error) {
	contractAddress := common.HexToAddress(config.CONTRACT_ADDRESS)
	zeroAddress := common.Address{}

	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	// Get latest block to calculate safe range (max 10000 blocks)
	latestBlock, err := rpc.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	fromBlock := int64(latestBlock) - 10000
	if fromBlock < 0 {
		fromBlock = 0
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(int64(latestBlock)),
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{transferSig}},
	}

	logs, err := rpc.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var transfers []models.TokenTransfer

	for _, log := range logs {
		if len(log.Topics) < 3 {
			continue
		}

		from := common.HexToAddress(log.Topics[1].Hex())
		to := common.HexToAddress(log.Topics[2].Hex())

		// Ignore mint events (Transfer from zero address) to focus on wallet-to-wallet transfers.
		if from == zeroAddress {
			continue
		}

		value := new(big.Int).SetBytes(log.Data)

		transfers = append(transfers, models.TokenTransfer{
			TxHash: log.TxHash.Hex(),
			From:   from.Hex(),
			To:     to.Hex(),
			Value:  value.String(),
		})
	}

	return transfers, nil
}
