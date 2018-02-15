package resp

import "github.com/MrDefinite/gedis/basicdata"

type protocolEncoder struct {
}

func createNewEncoder() *protocolEncoder {
	return &protocolEncoder{}
}

// The global used protocol encoder
var Encoder = createNewEncoder()

func (pe *protocolEncoder) EncodeToResponseObj(obj *basicdata.GedisObject) *basicdata.GedisObject {
	var res *basicdata.GedisObject
	t := basicdata.GetType(obj)
	switch t {
	case basicdata.GedisString:
		res = basicdata.AppendObject(CommonObjects.Plus, obj, CommonObjects.Crlf)
		break
	case basicdata.GedisList:
		break
	case basicdata.GedisHash:
		break
	case basicdata.GedisSet:
		break
	case basicdata.GedisZset:
		break
	}

	return res
}
