package main

import (
	"context"
	"explorer/api"
	"explorer/rpc"
	"explorer/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	rpc.InitClient()
	rpc.InitWsClient()

	go services.RunBlockHub()
	go services.RunTransferHub()
	go func() {
		if err := services.StartBlockStream(context.Background()); err != nil {
			log.Printf("block stream stopped: %v", err)
		}
	}()
	go func() {
		if err := services.StartTransferStream(context.Background()); err != nil {
			log.Printf("transfer stream stopped: %v", err)
		}
	}()

	r := gin.Default()
	r.GET("/config", api.ConfigHandler)
	r.GET("/block", api.BlockHandler)
	r.GET("/txs", api.TxHandler)
	r.GET("/transfers", api.TokenHandler)
	r.GET("/ws/blocks", api.BlockWsHandler)
	r.GET("/ws/transfers", api.TransferWsHandler)

	cfg, _ := services.GetConfig()
	log.Println("Explorer running at http://localhost:8080")
	log.Printf("  RPC: %s", cfg.RpcUrl)
	log.Printf("  Chain ID: %d", cfg.ChainId)
	log.Printf("  Token Contract: %s", cfg.ContractAddress)
	log.Fatal(r.Run(":8080"))
}
