package types

import (
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

var (
	log           = logrus.New()
	CommonObjects CommonObjectsStruct
)

const (
	GEDIS_STRING = 0
	GEDIS_LIST   = 1
	GEDIS_HASH   = 2
	GEDIS_SET    = 3
	GEDIS_ZSET   = 4
)

const (
	GEDIS_ENCODING_INT        = 0
	GEDIS_ENCODING_EMBSTR     = 1
	GEDIS_ENCODING_RAW        = 2
	GEDIS_ENCODING_HT         = 3
	GEDIS_ENCODING_LINKEDLIST = 4
	GEDIS_ENCODING_ZIPLIST    = 5
	GEDIS_ENCODING_INTSET     = 6
	GEDIS_ENCODING_SKIPLIST   = 7
)

type GedisObject struct {
	objType  uint8
	encoding uint8
	lru      uint32
	refCount int
	ptr      interface{}
}

type CommonObjectsStruct struct {
	Crlf, Ok, Err, Emptybulk, Czero, Cone, Cnegone, Pong, Space,
	Colon, Nullbulk, Nullmultibulk, Queued,
	Emptymultibulk, Wrongtypeerr, Nokeyerr, Syntaxerr, Sameobjecterr,
	Outofrangeerr, Noscripterr, Loadingerr, Slowscripterr, Bgsaveerr,
	Masterdownerr, Roslaveerr, Execaborterr, Noautherr, Noreplicaserr,
	Busykeyerr, Oomerr, Plus, messagebulk, pmessagebulk, subscribebulk,
	unsubscribebulk, psubscribebulk, punsubscribebulk, del, rpop, lpop,
	lpush, Emptyscan, minstring, maxstring *GedisObject
}

func TryObjectEncoding(obj *GedisObject) *GedisObject {
	return obj
}

func createObject(objType uint8, value interface{}) *GedisObject {
	obj := GedisObject{
		objType:  objType,
		encoding: GEDIS_ENCODING_RAW,
		ptr:      value,
		refCount: 1,
		lru:      0,
	}

	return &obj
}

func InitCommonObjects() {
	log.Info("Initializing gedis common used objects now")

	CommonObjects = CommonObjectsStruct{
		Crlf:           createObject(GEDIS_STRING, "\r\n"),
		Ok:             createObject(GEDIS_STRING, "+OK\r\n"),
		Err:            createObject(GEDIS_STRING, "-ERR\r\n"),
		Emptybulk:      createObject(GEDIS_STRING, "$0\r\n\r\n"),
		Czero:          createObject(GEDIS_STRING, ":0\r\n"),
		Cone:           createObject(GEDIS_STRING, ":1\r\n"),
		Cnegone:        createObject(GEDIS_STRING, ":-1\r\n"),
		Nullbulk:       createObject(GEDIS_STRING, "$-1\r\n"),
		Nullmultibulk:  createObject(GEDIS_STRING, "*-1\r\n"),
		Emptymultibulk: createObject(GEDIS_STRING, "*0\r\n"),
		Pong:           createObject(GEDIS_STRING, "+PONG\r\n"),
		Queued:         createObject(GEDIS_STRING, "+QUEUED\r\n"),
		Emptyscan:      createObject(GEDIS_STRING, "*2\r\n$1\r\n0\r\n*0\r\n"),

		Wrongtypeerr:  createObject(GEDIS_STRING, "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"),
		Nokeyerr:      createObject(GEDIS_STRING, "-ERR no such key\r\n"),
		Syntaxerr:     createObject(GEDIS_STRING, "-ERR syntax error\r\n"),
		Sameobjecterr: createObject(GEDIS_STRING, "-ERR source and destination objects are the same\r\n"),
		Outofrangeerr: createObject(GEDIS_STRING, "-ERR index out of range\r\n"),
		Noscripterr:   createObject(GEDIS_STRING, "-NOSCRIPT No matching script. Please use EVAL.\r\n"),
		Loadingerr:    createObject(GEDIS_STRING, "-LOADING Redis is loading the dataset in memory\r\n"),
		Slowscripterr: createObject(GEDIS_STRING, "-BUSY Redis is busy running a script. You can only call SCRIPT KILL or SHUTDOWN NOSAVE.\r\n"),
		Masterdownerr: createObject(GEDIS_STRING, "-MASTERDOWN Link with MASTER is down and slave-serve-stale-data is set to 'no'.\r\n"),
		Bgsaveerr:     createObject(GEDIS_STRING, "-MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.\r\n"),
		Roslaveerr:    createObject(GEDIS_STRING, "-READONLY You can't write against a read only slave.\r\n"),
		Noautherr:     createObject(GEDIS_STRING, "-NOAUTH Authentication required.\r\n"),
		Oomerr:        createObject(GEDIS_STRING, "-OOM command not allowed when used memory > 'maxmemory'.\r\n"),
		Execaborterr:  createObject(GEDIS_STRING, "-EXECABORT Transaction discarded because of previous errors.\r\n"),
		Noreplicaserr: createObject(GEDIS_STRING, "-NOREPLICAS Not enough good slaves to write.\r\n"),
		Busykeyerr:    createObject(GEDIS_STRING, "-BUSYKEY Target key name already exists.\r\n"),

		Space: createObject(GEDIS_STRING, " "),
		Colon: createObject(GEDIS_STRING, ":"),
		Plus:  createObject(GEDIS_STRING, "+"),
	}
}

func GetEncodeString(obj *GedisObject) (*GedisObject, error) {
	encodeObj := GedisObject{
		objType:  GEDIS_STRING,
		encoding: GEDIS_ENCODING_RAW,
	}
	encode := obj.encoding
	switch encode {
	case GEDIS_ENCODING_INT:
		encodeObj.ptr = "Integer"
	case GEDIS_ENCODING_EMBSTR:
		encodeObj.ptr = "Embed string"
	case GEDIS_ENCODING_RAW:
		encodeObj.ptr = "string"
	case GEDIS_ENCODING_HT:
		encodeObj.ptr = "Hash table"
	case GEDIS_ENCODING_LINKEDLIST:
		encodeObj.ptr = "Linked list"
	case GEDIS_ENCODING_ZIPLIST:
		encodeObj.ptr = "Zip list"
	case GEDIS_ENCODING_INTSET:
		encodeObj.ptr = "Int set"
	case GEDIS_ENCODING_SKIPLIST:
		encodeObj.ptr = "Skip list"
	default:
		return nil, errors.New("mal-format encode type found!")
	}
	return &encodeObj, nil
}
