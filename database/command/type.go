package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/protocol/resp"
)

type typeCommandProc struct {
}

func (c typeCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]

	if key == nil {
		return resp.CommonObjects.NullBulk
	}

	obj, err := basicdata.GetEncodeString(key)
	if err != nil {
		return resp.CommonObjects.Err
	}

	return obj
}
