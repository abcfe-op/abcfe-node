package storage

import (
	"fmt"

	log "github.com/abcfe/abcfe-node/common/logger"
	"github.com/abcfe/abcfe-node/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	db *leveldb.DB
}

func InitDB(cfg *config.Config) (*leveldb.DB, error) {
	dbName := fmt.Sprintf("leveldb_%d.db", cfg.Common.Port)
	dbPath := fmt.Sprintf("%s%s", cfg.DB.Path, dbName)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Error("Failed to load db: ", err)
	}

	return db, nil
}

func Close() {

}
