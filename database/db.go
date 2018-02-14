package database

import "github.com/MrDefinite/gedis/basicdata"

type GedisDB struct {
	// The key space for this DB
	Dict map[basicdata.GedisObject]interface{}

	// Timeout of keys with a timeout set
	Expires map[basicdata.GedisObject]interface{}

	// Keys with clients waiting for data (BLPOP)
	BlockingKeys map[basicdata.GedisObject]interface{}

	// Blocked keys that received a PUSH
	ReadyKeys map[basicdata.GedisObject]interface{}

	// WATCHED keys for MULTI/EXEC CAS
	WatchedKeys map[basicdata.GedisObject]interface{}

	// Database ID
	Id int

	// Average TTL, just for stats
	avgTTL int64
}

func InitDBFromZero() *GedisDB {
	db := GedisDB{
		Id:     0,
		avgTTL: -1,
	}
	db.Dict = make(map[basicdata.GedisObject]interface{})
	db.Expires = make(map[basicdata.GedisObject]interface{})
	db.BlockingKeys = make(map[basicdata.GedisObject]interface{})
	db.ReadyKeys = make(map[basicdata.GedisObject]interface{})
	db.WatchedKeys = make(map[basicdata.GedisObject]interface{})

	return &db
}

func InitDBFromFile() {

}
