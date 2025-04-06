package protocol

const (
	// 블록체인 설정 정보
	PrefixNetworkConfig = "net:config"

	// 메타데이터 관련 접두사
	PrefixLatestHeight    = "m:h" // 현재 블록 높이
	PrefixLatestBlockHash = "m:blk"
	PrefixStaker          = "m:staker" // 스테이커 리스트

	// 블록 관련 접두사
	PrefixBlock       = "blk:"       // 블록 해시 - 블록 데이터
	PrefixBlockHeight = "blk:h:"     // 블록 높이 - 블록 해시
	PrefixBlockTxIdx  = "blk:h:idx:" // 블록 해시 + 인덱스 - 인덱스에 해당되는 트랜잭션 해시

	// 트랜잭션 관련 접두사
	PrefixTx       = "tx:"      // 트랜잭션 해시 - 트랜잭션 데이터
	PrefixTxStatus = "tx:stat:" // 트랜잭션 해시 - 상태 값
	PrefixTxBlock  = "tx:blk:"  // 트랜잭션 해시 - 블록 해시

	// UTXO 관련 접두사
	PrefixUtxo = "utxo:" // UTXO 해시 - UTXO 데이터

	// 계정 관련 접두사
	PrefixAccount       = "acc:"         // 지갑 주소 - 상태값 (스테이킹 상태, 밸런스, POS 참여중이라면 합의 중인 상태 등)
	PrefixAccountTx     = "acc:tx:"      // 지갑 주소 - 트랜잭션 해시 값
	PrefixAccountTxTo   = "acc:tx:to:"   // 지갑 주소 - 수신자에 해당 되는 트랜잭션
	PrefixAccountTxFrom = "acc:tx:from:" // 지갑 주소 - 발신자에 해당 되는 트랜잭션
	PrefixAccountUtox   = "acc:utxo:"    // 지갑 주소 - UTXO 해시값

	// 컨센서스 관련 접두사
	PrefixStakerInfo = "acc:staker:" // 지갑 주소 - 스테이킹 정보
)
