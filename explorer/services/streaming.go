package services

import (
    "context"
    "encoding/json"
    "log"
    "math/big"
    "time"

    "explorer/config"
    "explorer/models"
    "explorer/rpc"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
)

var (
	clients = make(map[chan []byte]bool)
	transferClients = make(map[chan []byte]bool)

	register = make(chan chan[]byte)
	transferRegister = make(chan chan[]byte)

	unregister = make(chan chan[]byte)
	transferUnregister = make(chan chan[]byte)

	broadcast = make(chan []byte)
	transferBroadcast = make(chan []byte)
)

func RunBlockHub() {
	for {
		select {
		case ch :=  <- register:
			clients[ch] = true
		
		case ch := <- unregister:
			if _, ok := clients[ch]; ok {
				delete(clients, ch)
				close(ch)
			}

		case payload := <- broadcast:
			for ch := range clients {
				select {
				case ch <- payload:
				default:
					delete(clients, ch)
					close(ch)
				}
			}
		}
	}
}

func RunTransferHub() {
	for {
		select {
		case ch :=  <- transferRegister:
			transferClients[ch] = true
		
		case ch := <- transferUnregister:
			if _, ok := transferClients[ch]; ok {
				delete(transferClients, ch)
				close(ch)
			}

        case payload := <- transferBroadcast:
			for ch := range transferClients {
				select {
				case ch <- payload:
				default:
					delete(transferClients, ch)
					close(ch)
				}
			}
		}
	}
}

func RegisterBlockClient(ch chan []byte) {
	register <- ch 
}
func RegisterTransferClient(ch chan []byte) {
	transferRegister <- ch 
}

func UnregisterBlockClient(ch chan []byte) {
	unregister <- ch 
}
func UnregisterTransferClient(ch chan []byte) {
	transferUnregister <- ch 
}

func StartBlockStream(ctx context.Context) error {
    headers := make(chan *types.Header)
    sub, err := rpc.WsClient.SubscribeNewHead(ctx, headers)
    if err != nil {
        log.Printf("ws SubscribeNewHead failed, falling back to polling: %v", err)
        return startPollingBlockStream(ctx, 1500*time.Millisecond)
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case err := <-sub.Err():
            log.Printf("ws subscription dropped, falling back to polling: %v", err)
            return startPollingBlockStream(ctx, 1500*time.Millisecond)

        case header := <-headers:
            block := models.Block{
                Number:  header.Number.Uint64(),
                Gas:     header.GasUsed,
                Hash:    header.Hash().Hex(),
                TxCount: 0,
            }

            // Enrich via HTTP for tx count and canonical gas used.
            full, err := rpc.Client.BlockByHash(ctx, header.Hash())
            if err == nil && full != nil {
                block.TxCount = len(full.Transactions())
                block.Gas = full.GasUsed()
            }

            if err := broadcastBlock(block); err != nil {
                log.Printf("marshal block %d failed: %v", block.Number, err)
            }
        }
    }
}

func startPollingBlockStream(ctx context.Context, interval time.Duration) error {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastNumber uint64
    hasLast := false

    // Seed baseline so we only broadcast on actual changes.
    n, err := rpc.Client.BlockNumber(ctx)
    if err == nil {
        lastNumber = n
        hasLast = true
    } else {
        log.Printf("initial block number read failed: %v", err)
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case <-ticker.C:
            latest, err := rpc.Client.BlockNumber(ctx)
            if err != nil {
                log.Printf("poll block number failed: %v", err)
                continue
            }

            if !hasLast {
                lastNumber = latest
                hasLast = true
                continue
            }

            if latest == lastNumber {
                continue
            }

            // Broadcast each newly observed block in order.
            for n := lastNumber + 1; n <= latest; n++ {
                b, err := rpc.Client.BlockByNumber(ctx, new(big.Int).SetUint64(n))
                if err != nil {
                    log.Printf("poll block %d failed: %v", n, err)
                    break
                }

                mapped := models.Block{
                    Number:  b.Number().Uint64(),
                    Gas:     b.GasUsed(),
                    Hash:    b.Hash().Hex(),
                    TxCount: len(b.Transactions()),
                }

                if err := broadcastBlock(mapped); err != nil {
                    log.Printf("marshal block %d failed: %v", mapped.Number, err)
                }
            }

            lastNumber = latest
        }
    }
}

func broadcastBlock(block models.Block) error {
    payload, err := json.Marshal(block)
    if err != nil {
        return err
    }
    broadcast <- payload
    return nil
}

func StartTransferStream(ctx context.Context) error {
    headers := make(chan *types.Header)
    sub, err := rpc.WsClient.SubscribeNewHead(ctx, headers)
    if err != nil {
        log.Printf("ws transfer SubscribeNewHead failed, falling back to polling: %v", err)
        return startPollingTransferStream(ctx, 1500*time.Millisecond)
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case err := <-sub.Err():
            log.Printf("ws transfer subscription dropped, falling back to polling: %v", err)
            return startPollingTransferStream(ctx, 1500*time.Millisecond)

        case header := <-headers:
            transfers, err := getTransfersByBlock(ctx, header.Number.Uint64())
            if err != nil {
                log.Printf("load transfers for block %d failed: %v", header.Number.Uint64(), err)
                continue
            }

            for _, transfer := range transfers {
                if err := broadcastTransfer(transfer); err != nil {
                    log.Printf("marshal transfer in block %d failed: %v", header.Number.Uint64(), err)
                }
            }
        }
    }
}

func startPollingTransferStream(ctx context.Context, interval time.Duration) error {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    var lastNumber uint64
    hasLast := false

    n, err := rpc.Client.BlockNumber(ctx)
    if err == nil {
        lastNumber = n
        hasLast = true
    } else {
        log.Printf("initial transfer polling block number read failed: %v", err)
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case <-ticker.C:
            latest, err := rpc.Client.BlockNumber(ctx)
            if err != nil {
                log.Printf("poll transfer block number failed: %v", err)
                continue
            }

            if !hasLast {
                lastNumber = latest
                hasLast = true
                continue
            }

            if latest == lastNumber {
                continue
            }

            for n := lastNumber + 1; n <= latest; n++ {
                transfers, err := getTransfersByBlock(ctx, n)
                if err != nil {
                    log.Printf("poll transfers for block %d failed: %v", n, err)
                    continue
                }

                for _, transfer := range transfers {
                    if err := broadcastTransfer(transfer); err != nil {
                        log.Printf("marshal transfer in block %d failed: %v", n, err)
                    }
                }
            }

            lastNumber = latest
        }
    }
}

func getTransfersByBlock(ctx context.Context, blockNumber uint64) ([]models.TokenTransfer, error) {
    contractAddress := common.HexToAddress(config.CONTRACT_ADDRESS)
    zeroAddress := common.Address{}
    transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

    blockN := new(big.Int).SetUint64(blockNumber)
    query := ethereum.FilterQuery{
        FromBlock: blockN,
        ToBlock:   blockN,
        Addresses: []common.Address{contractAddress},
        Topics:    [][]common.Hash{{transferSig}},
    }

    logs, err := rpc.Client.FilterLogs(ctx, query)
    if err != nil {
        return nil, err
    }

    transfers := make([]models.TokenTransfer, 0, len(logs))
    for _, lg := range logs {
        if len(lg.Topics) < 3 {
            continue
        }

        from := common.HexToAddress(lg.Topics[1].Hex())
        to := common.HexToAddress(lg.Topics[2].Hex())

        if from == zeroAddress {
            continue
        }

        value := new(big.Int).SetBytes(lg.Data)
        transfers = append(transfers, models.TokenTransfer{
            TxHash: lg.TxHash.Hex(),
            From:   from.Hex(),
            To:     to.Hex(),
            Value:  value.String(),
        })
    }

    return transfers, nil
}

func broadcastTransfer(transfer models.TokenTransfer) error {
    payload, err := json.Marshal(transfer)
    if err != nil {
        return err
    }
    transferBroadcast <- payload
    return nil
}