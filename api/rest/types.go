package rest

// 일반적인 응답 구조체
type RestResp struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 블록체인 상태 응답
type BlockchainStatResp struct {
	Height    uint64 `json:"height"`
	BlockHash string `json:"blockHash"`
}

// 블록 응답
type BlockResp struct {
	Header       BlockHeaderResp `json:"header"`
	Transactions []TxResp        `json:"transactions"` // 트랜잭션 ID 목록
}

type BlockHeaderResp struct {
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash"`   // 이전 블록 해시
	Version    string `json:"version"`    // 블록체인 프로토콜 버전
	Height     uint64 `json:"height"`     // 블록 높이 (uint64로 변경)
	MerkleRoot string `json:"merkleRoot"` // 트랜잭션 머클 루트
	Timestamp  int64  `json:"timestamp"`  // 블록 생성 시간 (Unix 타임스탬프)
	// StateRoot  Hash   `json:"stateRoot"`  // 상태 머클 루트 (UTXO 또는 계정 상태)
}

// 트랜잭션 응답
type TxResp struct {
	ID        string        `json:"id"`
	Version   string        `json:"version"`
	Timestamp int64         `json:"timestamp"`
	Inputs    []interface{} `json:"inputs"`
	Outputs   []interface{} `json:"outputs"`
	Memo      string        `json:"memo"`
}

type SubmitTxReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
	Memo   string `json:"memo"`
	Data   []byte `json:"data"`
}
