package consensus

import (
	conf "github.com/abcfe/abcfe-node/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type Consensus struct {
	stop chan struct{}
	Conf conf.Config
	DB   *leveldb.DB // db내의 mutex는 복사되면 안됨
}
