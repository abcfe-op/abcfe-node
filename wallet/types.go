package wallet

type CipherParams struct {
	IV string `json:"iv"` // 초기화 벡터
}

type KDFParams struct {
	DkLen int    `json:"dklen"` // 도출된 키 길이
	N     int    `json:"n"`     // CPU/메모리 비용
	P     int    `json:"p"`     // 병렬화 매개변수
	R     int    `json:"r"`     // 블록 크기
	Salt  string `json:"salt"`  // 솔트
}

type Crypto struct {
	Cipher       string       `json:"cipher"`     // "aes-128-ctr"
	CipherText   string       `json:"ciphertext"` // 암호화된 개인키
	CipherParams CipherParams `json:"cipherparams"`
	KDF          string       `json:"kdf"` // "scrypt"
	KDFParams    KDFParams    `json:"kdfparams"`
	MAC          string       `json:"mac"` // 무결성 검증
}
