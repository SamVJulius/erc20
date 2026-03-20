package services

import (
	"context"
	"explorer/config"
	"explorer/models"
	"explorer/rpc"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

var TransferEventSig = common.HexToHash("0xddf252ad")

func GetTokenTransfers() ([]models.TokenTransfer, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(config.CONTRACT_ADDRESS),
		},
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

		from := common.HexToAddress(string(log.Topics[1].Hex()))
		to := common.HexToAddress(string(log.Topics[2].Hex()))

		value := new(big.Int).SetBytes(log.Data)

		transfers = append(transfers, models.TokenTransfer{
			TxHash: log.TxHash.Hex(),
			From: from.Hex(),
			To: to.Hex(),
			Value: value.String(),
		})
	}

	return transfers, nil
}