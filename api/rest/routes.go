package rest

import (
	"net/http"

	"github.com/abcfe/abcfe-node/core"
	"github.com/gorilla/mux"
)

func setupRouter(blockchain *core.BlockChain) http.Handler {
	r := mux.NewRouter()

	// 미들웨어 설정
	r.Use(LoggingMiddleware)
	r.Use(RecoveryMiddleware)

	// 기본 경로
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// 블록체인 API 경로
	api := r.PathPrefix("/api/v1").Subrouter()

	// 블록체인 상태
	api.HandleFunc("/status", GetStatus(blockchain)).Methods("GET")

	// 블록 관련 API
	// api.HandleFunc("/block", ComposeAndAddBlock(blockchain)).Methods("POST") // 테스트 전용 블록 구성 및 블록 추가 (검증은 없음)
	api.HandleFunc("/block/latest", GetLatestBlock(blockchain)).Methods("GET")
	api.HandleFunc("/block/{height}", GetBlockByHeight(blockchain)).Methods("GET")
	api.HandleFunc("/block/hash/{hash}", GetBlockByHash(blockchain)).Methods("GET")

	// 트랜잭션 관련 API
	// api.HandleFunc("/tx", SubmitTx(blockchain)).Methods("POST") // tx to mempool
	api.HandleFunc("/tx/{txid}", GetTx(blockchain)).Methods("GET")

	// UTXO 관련 API
	api.HandleFunc("/address/{address}/utxo", GetAddressUtxo(blockchain)).Methods("GET")
	api.HandleFunc("/address/{address}/balance", GetBalanceByUtxo(blockchain)).Methods("GET")

	// 계정 관련
	// api.HandleFunc("/address/{address}", GetAddress(blockchain)).Methods("GET")

	return r
}
