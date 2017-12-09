package types

const (
	GEDIS_STRING = 0
	GEDIS_LIST
	GEDIS_HASH
	GEDIS_SET
	GEDIS_ZSET
)


type gedisObject struct {
	objType  uint8
	encoding uint8
	lru      uint32
	refCount int
}
