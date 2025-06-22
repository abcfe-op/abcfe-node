package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/abcfe/abcfe-node/common/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/db_browser.go <db_path> [command]")
		fmt.Println("Commands:")
		fmt.Println("  meta     - Show metadata")
		fmt.Println("  blocks   - List all blocks")
		fmt.Println("  txs      - List all transactions")
		fmt.Println("  block <height> - Show specific block")
		fmt.Println("  tx <hash> - Show specific transaction")
		fmt.Println("  all      - Show all data")
		return
	}

	dbPath := os.Args[1]
	command := "meta"
	if len(os.Args) > 2 {
		command = os.Args[2]
	}

	// LevelDB 열기
	db, err := leveldb.OpenFile(dbPath, &opt.Options{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fmt.Printf("Database opened: %s\n\n", dbPath)

	switch command {
	case "meta":
		showMetadata(db)
	case "blocks":
		listBlocks(db)
	case "txs":
		listTransactions(db)
	case "block":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run tools/db_browser.go <db_path> block <height>")
			return
		}
		height, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			fmt.Printf("Invalid height: %v\n", err)
			return
		}
		showBlock(db, height)
	case "tx":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run tools/db_browser.go <db_path> tx <hash>")
			return
		}
		showTransaction(db, os.Args[3])
	case "all":
		showAllData(db)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func showMetadata(db *leveldb.DB) {
	fmt.Println("=== METADATA ===")

	// 최신 높이
	heightBytes, err := db.Get([]byte("meta:height"), nil)
	if err != nil {
		fmt.Printf("Latest Height: Not found (%v)\n", err)
	} else {
		fmt.Printf("Latest Height: %s\n", string(heightBytes))
	}

	// 최신 블록 해시
	hashBytes, err := db.Get([]byte("meta:hash"), nil)
	if err != nil {
		fmt.Printf("Latest Block Hash: Not found (%v)\n", err)
	} else {
		fmt.Printf("Latest Block Hash: %s\n", string(hashBytes))
	}

	fmt.Println()
}

func listBlocks(db *leveldb.DB) {
	fmt.Println("=== BLOCKS ===")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	prefix := []byte("blk:h:")
	for iter.Seek(prefix); iter.Valid() && iter.Key()[0] == prefix[0]; iter.Next() {
		key := string(iter.Key())
		if len(key) >= len("blk:h:") && key[:len("blk:h:")] == "blk:h:" {
			height := key[len("blk:h:"):]
			hash := string(iter.Value())
			fmt.Printf("Height %s: %s\n", height, hash)
		}
	}
	fmt.Println()
}

func listTransactions(db *leveldb.DB) {
	fmt.Println("=== TRANSACTIONS ===")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	prefix := []byte("tx:")
	count := 0
	for iter.Seek(prefix); iter.Valid() && iter.Key()[0] == prefix[0]; iter.Next() {
		key := string(iter.Key())
		if len(key) >= len("tx:") && key[:len("tx:")] == "tx:" && !contains(key, ":") {
			// tx:로 시작하고 추가 콜론이 없는 경우만 (전체 트랜잭션 데이터)
			txHash := key[len("tx:"):]
			fmt.Printf("Transaction: %s\n", txHash)
			count++
		}
	}
	fmt.Printf("Total transactions: %d\n\n", count)
}

func showBlock(db *leveldb.DB, height uint64) {
	fmt.Printf("=== BLOCK %d ===\n", height)

	// 높이로 블록 해시 조회
	heightKey := utils.GetBlockHeightKey(height)
	blockHashBytes, err := db.Get(heightKey, nil)
	if err != nil {
		fmt.Printf("Block not found: %v\n", err)
		return
	}

	blockHashStr := string(blockHashBytes)
	fmt.Printf("Block Hash: %s\n", blockHashStr)

	// 블록 해시로 블록 데이터 조회
	blockHash, err := utils.StringToHash(blockHashStr)
	if err != nil {
		fmt.Printf("Invalid block hash: %v\n", err)
		return
	}

	blockKey := utils.GetBlockHashKey(blockHash)
	blockData, err := db.Get(blockKey, nil)
	if err != nil {
		fmt.Printf("Block data not found: %v\n", err)
		return
	}

	fmt.Printf("Block Data Size: %d bytes\n", len(blockData))
	fmt.Printf("Block Data (hex): %s\n", hex.EncodeToString(blockData[:min(100, len(blockData))]))
	if len(blockData) > 100 {
		fmt.Println("... (truncated)")
	}
	fmt.Println()
}

func showTransaction(db *leveldb.DB, txHashStr string) {
	fmt.Printf("=== TRANSACTION %s ===\n", txHashStr)

	// 트랜잭션 해시 변환
	txHash, err := utils.StringToHash(txHashStr)
	if err != nil {
		fmt.Printf("Invalid transaction hash: %v\n", err)
		return
	}

	// 트랜잭션 데이터 조회
	txKey := utils.GetTxHashKey(txHash)
	txData, err := db.Get(txKey, nil)
	if err != nil {
		fmt.Printf("Transaction not found: %v\n", err)
		return
	}

	fmt.Printf("Transaction Data Size: %d bytes\n", len(txData))
	fmt.Printf("Transaction Data (hex): %s\n", hex.EncodeToString(txData[:min(100, len(txData))]))
	if len(txData) > 100 {
		fmt.Println("... (truncated)")
	}

	// 트랜잭션이 포함된 블록 해시 조회
	txBlockKey := utils.GetTxBlockHashKey(txHash)
	blockHashBytes, err := db.Get(txBlockKey, nil)
	if err == nil {
		blockHash := utils.BytesToHash(blockHashBytes)
		fmt.Printf("Included in Block: %s\n", utils.HashToString(blockHash))
	}

	fmt.Println()
}

func showAllData(db *leveldb.DB) {
	fmt.Println("=== ALL DATABASE DATA ===")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	count := 0
	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		fmt.Printf("[%d] Key: %s\n", count, key)
		fmt.Printf("     Value Size: %d bytes\n", len(value))
		if len(value) <= 100 {
			fmt.Printf("     Value: %s\n", string(value))
		} else {
			fmt.Printf("     Value (hex): %s...\n", hex.EncodeToString(value[:50]))
		}
		fmt.Println()

		count++
		if count >= 50 { // 최대 50개만 표시
			fmt.Printf("... (showing first 50 entries)\n")
			break
		}
	}

	fmt.Printf("Total entries: %d\n", count)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
