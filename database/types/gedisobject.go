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
	Crlf, Ok, Err, Emptybulk, Czero, Cone, cnegone, pong, space,
	colon, Nullbulk, nullmultibulk, queued,
	emptymultibulk, wrongtypeerr, nokeyerr, syntaxerr, sameobjecterr,
	outofrangeerr, noscripterr, loadingerr, slowscripterr, bgsaveerr,
	masterdownerr, roslaveerr, execaborterr, noautherr, noreplicaserr,
	busykeyerr, oomerr, plus, messagebulk, pmessagebulk, subscribebulk,
	unsubscribebulk, psubscribebulk, punsubscribebulk, del, rpop, lpop,
	lpush, emptyscan, minstring, maxstring *GedisObject
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
		cnegone:        createObject(GEDIS_STRING, ":-1\r\n"),
		Nullbulk:       createObject(GEDIS_STRING, "$-1\r\n"),
		nullmultibulk:  createObject(GEDIS_STRING, "*-1\r\n"),
		emptymultibulk: createObject(GEDIS_STRING, "*0\r\n"),
		pong:           createObject(GEDIS_STRING, "+PONG\r\n"),
		queued:         createObject(GEDIS_STRING, "+QUEUED\r\n"),
		emptyscan:      createObject(GEDIS_STRING, "*2\r\n$1\r\n0\r\n*0\r\n"),
	}
}

func GetEncodeString(obj *GedisObject) (*GedisObject, error) {
	encodeObj := GedisObject{
		objType: GEDIS_STRING,
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
