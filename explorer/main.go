package main

import (
	"explorer/api"
	"explorer/rpc"
	"explorer/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	rpc.InitClient()
	rpc.InitWsClient()

	r := gin.Default()
	r.GET("/config", api.ConfigHandler)
	r.GET("/block", api.BlockHandler)
	r.GET("/txs", api.TxHandler)
	r.GET("/transfers", api.TokenHandler)

	cfg, _ := services.GetConfig()
	log.Println("Explorer running at http://localhost:8080")
	log.Printf("  RPC: %s", cfg.RpcUrl)
	log.Printf("  Chain ID: %d", cfg.ChainId)
	log.Printf("  Token Contract: %s", cfg.ContractAddress)
	log.Fatal(r.Run(":8080"))
}
