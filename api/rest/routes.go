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
	api.HandleFunc("/status", StatusHandler(blockchain)).Methods("GET")

	// 블록 관련 API
	api.HandleFunc("/block/latest", GetLatestBlockHandler(blockchain)).Methods("GET")
	api.HandleFunc("/block/{height}", GetBlockByHeightHandler(blockchain)).Methods("GET")
	api.HandleFunc("/block/hash/{hash}", GetBlockByHashHandler(blockchain)).Methods("GET")

	// api.HandleFunc("/block", AddBlockHandler(blockchain)).Methods("POST")

	// 트랜잭션 관련 API
	api.HandleFunc("/tx/{txid}", GetTxHandler(blockchain)).Methods("GET")
	// api.HandleFunc("/tx", SubmitTxHandler(blockchain)).Methods("POST")

	// UTXO 관련 API
	// api.HandleFunc("/address/{address}/utxo", GetAddressUTXOHandler(blockchain)).Methods("GET")
	// api.HandleFunc("/address/{address}/balance", GetAddressBalanceHandler(blockchain)).Methods("GET")

	// 계정 관련
	// api.HandleFunc("/account/{address}", GetAccountHandler(blockchain)).Methods("GET")
	// api.HandleFunc("/account/{address}/balance", GetAccountBalanceHandler(blockchain)).Methods("GET")

	return r
}
