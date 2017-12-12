package database

import "github.com/MrDefinite/gedis/database/types"

type GedisDB struct {
	// The key space for this DB
	Dict map[types.GedisObject]interface{}

	// Timeout of keys with a timeout set
	Expires map[types.GedisObject]interface{}

	// Keys with clients waiting for data (BLPOP)
	BlockingKeys map[types.GedisObject]interface{}

	// Blocked keys that received a PUSH
	ReadyKeys map[types.GedisObject]interface{}

	// WATCHED keys for MULTI/EXEC CAS
	WatchedKeys map[types.GedisObject]interface{}

	// Database ID
	Id int

	// Average TTL, just for stats
	avgTTL int64
}

func InitDBFromZero() *GedisDB {
	db := GedisDB{
		Id: 0,
		avgTTL: -1,
	}
	db.Dict = make(map[types.GedisObject]interface{})
	db.Expires = make(map[types.GedisObject]interface{})
	db.BlockingKeys = make(map[types.GedisObject]interface{})
	db.ReadyKeys = make(map[types.GedisObject]interface{})
	db.WatchedKeys = make(map[types.GedisObject]interface{})

	return &db
}

func InitDBFromFile()  {

}
