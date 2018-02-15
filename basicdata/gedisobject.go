package basicdata

import (
	"github.com/pkg/errors"
)

var (
	ErrEncodingType = errors.New("mal-format encode type found!")
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

func TryObjectEncoding(obj *GedisObject) *GedisObject {
	return obj
}

func CreateObject(objType GedisObjectType, value interface{}) *GedisObject {
	obj := GedisObject{
		objType:  objType,
		encoding: GedisEncodingRaw,
		ptr:      value,
		refCount: 1,
		lru:      0,
	}

	return &obj
}

func CreateObjectWithEncoding(objType GedisObjectType, encoding GedisObjectEncodingType, value interface{}) *GedisObject {
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

func AppendObject(appended ...*GedisObject) *GedisObject {
	// TODO do some analyse here
	res := CreateEmptyStringObject()

	var val string
	for _, op := range appended {
		switch p := op.ptr.(type) {
		case string:
			val += p
		}
	}
	res.ptr = val
	return res
}
