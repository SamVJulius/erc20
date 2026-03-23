package api

import (
	"explorer/services"

	"github.com/gin-gonic/gin"
)

func BlockHandler(c *gin.Context) {
	block, err := services.GetLatestBlock()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, block)
}

func TxHandler(c *gin.Context) {
	txs, err := services.GetLatestTransactions()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, txs)
}

func TokenHandler(c *gin.Context) {
	transfers, err := services.GetContractTransfers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, transfers)
}

func ConfigHandler(c *gin.Context) {
	cfg, err := services.GetConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, cfg)
}
