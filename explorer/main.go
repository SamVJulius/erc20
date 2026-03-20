package main

import (
	"explorer/api"
	"explorer/rpc"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	rpc.InitClient()

	r := gin.Default()
	r.GET("/block", api.BlockHandler)
	r.GET("/txs", api.TxHandler)
	r.GET("/transfers", api.TokenHandler)

	log.Println("Explorer running at http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
