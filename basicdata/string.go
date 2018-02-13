package basicdata

import (
	"github.com/pkg/errors"
)

var (
	ErrNonString         = errors.New("Cannot get string from non string object")
	ErrNonStringEncoding = errors.New("Cannot get string from non string encoding object")
)

func CreateStringObjectWithBytes(in []byte) *GedisObject {
	t := getDesiredEncodingForBytes(in)
	return createObjectWithEncoding(GedisString, t, string(in))
}

func CreateStringObject(in string) *GedisObject {
	t := getDesiredEncoding(in)
	return createObjectWithEncoding(GedisString, t, in)
}

func getDesiredEncoding(in string) GedisObjectEncodingType {
	// TODO: detect data length
	return GedisEncodingRaw
}

func getDesiredEncodingForBytes(in []byte) GedisObjectEncodingType {
	// TODO: detect data length
	return GedisEncodingRaw
}

func GetStringValueFromObject(obj *GedisObject) (string, error) {
	if obj.objType != GedisString {
		return "", ErrNonString
	}

	if obj.encoding == GedisEncodingRaw {
		data, err := getRawStringValueFromObject(obj.ptr)
		if err != nil {
			return "", err
		}
		return data, nil
	}
	if obj.encoding == GedisEncodingEmbstr {
		data, err := getEmbStrStringValueFromObject(obj.ptr)
		if err != nil {
			return "", err
		}
		return data, nil
	}

	return "", ErrNonString
}

func getRawStringValueFromObject(data interface{}) (string, error) {
	switch _ := data.(type) {
	case string:
		return data.(string), nil
	}
	return "", ErrNonStringEncoding
}

func getEmbStrStringValueFromObject(data interface{}) (string, error) {
	// TODO: change it later
	switch _ := data.(type) {
	case string:
		return data.(string), nil
	}
	return "", ErrNonStringEncoding
}
