package types

import "fmt"

func CreateStringObject(in string) *GedisObject {
	// TODO: use GEDIS_ENCODING_EMBSTR or GEDIS_ENCODING_RAW ?
	return createObject(GEDIS_STRING, in)
}

func GetStringValueFromObject(obj *GedisObject) string {
	if obj.objType != GEDIS_STRING {
		fmt.Errorf("this is not a string object, it is encoded with %s", obj.encoding)
	}
	return obj.ptr.(string)
}

