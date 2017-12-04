package database


type GedisDB struct {
	// The key space for this DB
	dict *map[string]interface{}

	// Timeout of keys with a timeout set
	expires *map[string]interface{}

	// Keys with clients waiting for data (BLPOP)
	blockingKeys *map[string]interface{}

	// Blocked keys that received a PUSH
	readyKeys *map[string]interface{}

	// WATCHED keys for MULTI/EXEC CAS
	watchedKeys *map[string]interface{}

	// Database ID
	id int

	// Average TTL, just for stats
	avgTTL int64
}

