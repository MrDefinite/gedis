package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
)

type typeCommandProc struct {
}

func (c typeCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]

	if key == nil {
		return basicdata.CommonObjects.NullBulk
	}

	obj, err := basicdata.GetEncodeString(key)
	if err != nil {
		return basicdata.CommonObjects.Err
	}

	return obj
}
