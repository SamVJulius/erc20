package explorer

import (
	"explorer/api"
	"explorer/rpc"
	"log"
	"net/http"
)

func main() {
	rpc.InitClient()

	http.HandleFunc("/block", api.BlockHandler)
	http.HandleFunc("/txs", api.TxHandler)
	http.HandleFunc("/transfers", api.TokenHandler)

	log.Println("Explorer running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}