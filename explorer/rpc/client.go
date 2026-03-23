package rpc

import (
	"log"
	"explorer/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

var Client *ethclient.Client
var WsClient *ethclient.Client

func InitClient() {
	client, err := ethclient.Dial(config.RPC_URL)
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	Client = client
}

func InitWsClient() {
	wsclient, err := ethclient.Dial(config.RPC_WS_URL)
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	WsClient = wsclient
}