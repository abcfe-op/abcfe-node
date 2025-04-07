package core

// 트랜잭션 구조체
type Transaction struct {
	Version   uint32 `json:"version"`   // 트랜잭션 버전
	ID        Hash   `json:"id"`        // 트랜잭션 ID (해시)
	Timestamp int64  `json:"timestamp"` // 트랜잭션 생성 시간

	Inputs  []*TxInput  `json:"inputs"`  // 트랜잭션 입력
	Outputs []*TxOutput `json:"outputs"` // 트랜잭션 출력

	Memo string `json:"memo"` // 트랜잭션 메모 (inputData 대체)
	Data []byte `json:"data"` // 임의 데이터 (스마트 컨트랙트 호출 등)
}

type TxInput struct {
	TxID        Hash      `json:"txId"`        // 참조 트랜잭션 ID
	OutputIndex uint32    `json:"outputIndex"` // 참조 출력 인덱스
	Signature   Signature `json:"signature"`   // 서명
	PublicKey   []byte    `json:"publicKey"`   // 공개키
	Sequence    uint32    `json:"sequence"`    // 시퀀스 번호 (RBF 지원)
}

type TxOutput struct {
	Address Address `json:"address"` // 수신자 주소
	Amount  uint64  `json:"amount"`  // 금액 (int에서 uint64로 변경)
	TxType  uint8   `json:"txType"`  // 스크립트 타입 (일반/스테이킹/기타)
}
