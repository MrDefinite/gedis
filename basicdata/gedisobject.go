package basicdata

import (
	"github.com/pkg/errors"
)

var (
	ErrEncodingType = errors.New("mal-format encode type found!")
)

var (
	CommonObjects = CommonObjectStruct{
		Crlf:           createObject(GedisString, "\r\n"),
		Ok:             createObject(GedisString, "+OK\r\n"),
		Err:            createObject(GedisString, "-ERR\r\n"),
		EmptyBulk:      createObject(GedisString, "$0\r\n\r\n"),
		Zero:           createObject(GedisString, ":0\r\n"),
		One:            createObject(GedisString, ":1\r\n"),
		NegOne:         createObject(GedisString, ":-1\r\n"),
		NullBulk:       createObject(GedisString, "$-1\r\n"),
		NullMultiBulk:  createObject(GedisString, "*-1\r\n"),
		EmptyMultiBulk: createObject(GedisString, "*0\r\n"),
		Pong:           createObject(GedisString, "+PONG\r\n"),
		Queued:         createObject(GedisString, "+QUEUED\r\n"),
		EmptyScan:      createObject(GedisString, "*2\r\n$1\r\n0\r\n*0\r\n"),

		WrongTypeErr:  createObject(GedisString, "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"),
		NoKeyErr:      createObject(GedisString, "-ERR no such key\r\n"),
		SynTaxErr:     createObject(GedisString, "-ERR syntax error\r\n"),
		SameObjectErr: createObject(GedisString, "-ERR source and destination objects are the same\r\n"),
		OutOfRangeErr: createObject(GedisString, "-ERR index out of range\r\n"),
		NoScriptErr:   createObject(GedisString, "-NOSCRIPT No matching script. Please use EVAL.\r\n"),
		LoadingErr:    createObject(GedisString, "-LOADING Redis is loading the dataset in memory\r\n"),
		SlowScriptErr: createObject(GedisString, "-BUSY Redis is busy running a script. You can only call SCRIPT KILL or SHUTDOWN NOSAVE.\r\n"),
		MasterDownErr: createObject(GedisString, "-MASTERDOWN Link with MASTER is down and slave-serve-stale-data is set to 'no'.\r\n"),
		BgSaveErr:     createObject(GedisString, "-MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.\r\n"),
		RoSlaveErr:    createObject(GedisString, "-READONLY You can't write against a read only slave.\r\n"),
		NoAuthErr:     createObject(GedisString, "-NOAUTH Authentication required.\r\n"),
		OomErr:        createObject(GedisString, "-OOM command not allowed when used memory > 'maxmemory'.\r\n"),
		ExecAbortErr:  createObject(GedisString, "-EXECABORT Transaction discarded because of previous errors.\r\n"),
		NoReplicasErr: createObject(GedisString, "-NOREPLICAS Not enough good slaves to write.\r\n"),
		BusyKeyErr:    createObject(GedisString, "-BUSYKEY Target key name already exists.\r\n"),

		Space: createObject(GedisString, " "),
		Colon: createObject(GedisString, ":"),
		Plus:  createObject(GedisString, "+"),
	}
)

const (
	GedisString = 0
	GedisList   = 1
	GedisHash   = 2
	GedisSet    = 3
	GedisZset   = 4
)

type GedisObjectType uint8

const (
	GedisEncodingInt        = 0
	GedisEncodingEmbstr     = 1
	GedisEncodingRaw        = 2
	GedisEncodingHt         = 3
	GedisEncodingLinkedlist = 4
	GedisEncodingZiplist    = 5
	GedisEncodingIntset     = 6
	GedisEncodingSkiplist   = 7
)

type GedisObjectEncodingType uint8

type GedisObject struct {
	objType  GedisObjectType
	encoding GedisObjectEncodingType
	lru      uint32
	refCount int
	ptr      interface{}
}

type CommonObjectStruct struct {
	Crlf, Ok, Err, EmptyBulk, Zero, One, NegOne, Pong, Space,
	Colon, NullBulk, NullMultiBulk, Queued,
	EmptyMultiBulk, WrongTypeErr, NoKeyErr, SynTaxErr, SameObjectErr,
	OutOfRangeErr, NoScriptErr, LoadingErr, SlowScriptErr, BgSaveErr,
	MasterDownErr, RoSlaveErr, ExecAbortErr, NoAuthErr, NoReplicasErr,
	BusyKeyErr, OomErr, Plus, EmptyScan *GedisObject
}

func TryObjectEncoding(obj *GedisObject) *GedisObject {
	return obj
}

func createObject(objType GedisObjectType, value interface{}) *GedisObject {
	obj := GedisObject{
		objType:  objType,
		encoding: GedisEncodingRaw,
		ptr:      value,
		refCount: 1,
		lru:      0,
	}

	return &obj
}

func createObjectWithEncoding(objType GedisObjectType, encoding GedisObjectEncodingType, value interface{}) *GedisObject {
	obj := GedisObject{
		objType:  objType,
		encoding: encoding,
		ptr:      value,
		refCount: 1,
		lru:      0,
	}

	return &obj
}

func GetType(obj *GedisObject) GedisObjectType {
	return obj.objType
}

func GetEncodeString(obj *GedisObject) (*GedisObject, error) {
	encodeObj := GedisObject{
		objType:  GedisString,
		encoding: GedisEncodingRaw,
	}
	encode := obj.encoding
	switch encode {
	case GedisEncodingInt:
		encodeObj.ptr = "Integer"
	case GedisEncodingEmbstr:
		encodeObj.ptr = "Embed string"
	case GedisEncodingRaw:
		encodeObj.ptr = "string"
	case GedisEncodingHt:
		encodeObj.ptr = "Hash table"
	case GedisEncodingLinkedlist:
		encodeObj.ptr = "Linked list"
	case GedisEncodingZiplist:
		encodeObj.ptr = "Zip list"
	case GedisEncodingIntset:
		encodeObj.ptr = "Int set"
	case GedisEncodingSkiplist:
		encodeObj.ptr = "Skip list"
	default:
		return nil, ErrEncodingType
	}
	return &encodeObj, nil
}
