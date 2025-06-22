package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abcfe/abcfe-node/common/utils"
	"github.com/abcfe/abcfe-node/core"
	"github.com/gorilla/mux"
)

// send response
func sendResp(w http.ResponseWriter, statusCode int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := RestResp{
		Success: err == nil,
		Data:    data,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// get home response
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]string{
		"name":    "ABCFE Blockchain API",
		"version": "1.0.0",
	}
	sendResp(w, http.StatusOK, info, nil)
}

// get chain status response
func StatusHandler(bc *core.BlockChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := bc.GetChainStatus()

		response := BlockchainStatResp{
			Height:    status.LatestHeight,
			BlockHash: status.LatestBlockHash,
		}

		sendResp(w, http.StatusOK, response, nil)
	}
}

// get latest block response
func GetLatestBlockHandler(bc *core.BlockChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height, err := bc.GetLatestHeight()
		if err != nil {
			sendResp(w, http.StatusInternalServerError, nil, err)
			return
		}

		block, err := bc.GetBlockByHeight(height)
		if err != nil {
			sendResp(w, http.StatusInternalServerError, nil, err)
			return
		}

		response, err := formatBlockResp(block)
		if err != nil {
			sendResp(w, http.StatusInternalServerError, nil, err)
			return
		}

		sendResp(w, http.StatusOK, response, nil)
	}
}

// get block by height response
func GetBlockByHeightHandler(bc *core.BlockChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		heightStr := vars["height"]

		height, err := strconv.ParseUint(heightStr, 10, 64)
		if err != nil {
			sendResp(w, http.StatusBadRequest, nil, err)
			return
		}

		block, err := bc.GetBlockByHeight(height)
		if err != nil {
			sendResp(w, http.StatusNotFound, nil, err)
			return
		}

		response, err := formatBlockResp(block)
		if err != nil {
			sendResp(w, http.StatusInternalServerError, nil, err)
			return
		}

		sendResp(w, http.StatusOK, response, nil)
	}
}

// get block by hash response
func GetBlockByHashHandler(bc *core.BlockChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashStr := vars["hash"]

		hash, err := utils.StringToHash(hashStr)
		if err != nil {
			sendResp(w, http.StatusBadRequest, nil, err)
			return
		}

		block, err := bc.GetBlockByHash(hash)
		if err != nil {
			sendResp(w, http.StatusNotFound, nil, err)
			return
		}

		response, err := formatBlockResp(block)
		if err != nil {
			sendResp(w, http.StatusInternalServerError, nil, err)
			return
		}

		sendResp(w, http.StatusOK, response, nil)
	}
}

// get tx response
func GetTxHandler(bc *core.BlockChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		txIDStr := vars["txid"]

		txID, err := utils.StringToHash(txIDStr)
		if err != nil {
			sendResp(w, http.StatusBadRequest, nil, err)
			return
		}

		tx, err := bc.GetTx(txID)
		if err != nil {
			sendResp(w, http.StatusNotFound, nil, err)
			return
		}

		response := formatTxResp(tx)
		sendResp(w, http.StatusOK, response, nil)
	}
}

// get block response
func formatBlockResp(block *core.Block) (BlockResp, error) {
	txDetails := make([]TxResp, len(block.Transactions))
	for i, tx := range block.Transactions {
		txDetails[i] = formatTxResp(tx)
	}

	response := BlockResp{
		Header: BlockHeaderResp{
			Version:    block.Header.Version,
			Height:     block.Header.Height,
			PrevHash:   utils.HashToString(block.Header.PrevHash),
			MerkleRoot: utils.HashToString(block.Header.MerkleRoot),
			Timestamp:  block.Header.Timestamp,
		},
		Hash:         utils.HashToString(block.Hash),
		Transactions: txDetails,
	}

	return response, nil
}

// get tx response
func formatTxResp(tx *core.Transaction) TxResp {
	return TxResp{
		ID:        utils.HashToString(tx.ID),
		Version:   tx.Version,
		Timestamp: tx.Timestamp,
		Inputs:    formatTxInputs(tx.Inputs),
		Outputs:   formatTxOutputs(tx.Outputs),
		Memo:      tx.Memo,
	}
}

// get tx input response
func formatTxInputs(inputs []*core.TxInput) []interface{} {
	result := make([]interface{}, len(inputs))
	for i, input := range inputs {
		result[i] = map[string]interface{}{
			"txid":        utils.HashToString(input.TxID),
			"outputIndex": input.OutputIndex,
			"signature":   input.Signature,
			"publicKey":   input.PublicKey,
		}
	}
	return result
}

// get tx output response
func formatTxOutputs(outputs []*core.TxOutput) []interface{} {
	result := make([]interface{}, len(outputs))
	for i, output := range outputs {
		result[i] = map[string]interface{}{
			"address": utils.AddressToString(output.Address),
			"amount":  output.Amount,
			"txType":  output.TxType,
		}
	}
	return result
}
