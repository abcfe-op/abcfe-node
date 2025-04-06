package consensus

import (
	"encoding/json"

	"github.com/abcfe/abcfe-node/protocol"
	"github.com/syndtr/goleveldb/leveldb"
)

type Staker struct {
	Address string
	Amount  uint64
	AtStart uint64 // unix
	// TODO etc
}

func GetStakers(db *leveldb.DB) ([]string, error) {
	// db
	key := []byte(protocol.PrefixStaker)
	data, err := db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var stakerAddrs []string
	if err := json.Unmarshal(data, &stakerAddrs); err != nil {
		return nil, err
	}

	return stakerAddrs, nil
}
