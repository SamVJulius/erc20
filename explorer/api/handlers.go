package api

import (
	"explorer/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var blockWsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func BlockWsHandler(c *gin.Context) {
	conn, err := blockWsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	latestBlock, err := services.GetLatestBlock()
	if err != nil {
		_ = conn.Close()
		return
	}
	if err := conn.WriteJSON(latestBlock); err != nil {
		_ = conn.Close()
		return
	}

	clientCh := make(chan []byte, 64)
	services.RegisterBlockClient(clientCh)

	defer func() {
		services.UnregisterBlockClient(clientCh)
		_ = conn.Close()
	}()

	// Writer loop: push broadcast payloads to websocket client.
	go func() {
		for payload := range clientCh {
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	}()

	// Read loop: keep connection alive and detect disconnect.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func TransferWsHandler(c *gin.Context) {
	conn, err := blockWsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	transfers, err := services.GetContractTransfers()
	if err != nil {
		_ = conn.Close()
		return
	}
	if err := conn.WriteJSON(transfers); err != nil {
		_ = conn.Close()
		return
	}

	clientCh := make(chan []byte, 64)
	services.RegisterTransferClient(clientCh)

	defer func() {
		services.UnregisterTransferClient(clientCh)
		_ = conn.Close()
	}()

	// Writer loop: push broadcast payloads to websocket client.
	go func() {
		for payload := range clientCh {
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	}()

	// Read loop: keep connection alive and detect disconnect.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}

}