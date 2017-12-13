package command

import (
	"github.com/MrDefinite/gedis/database/types"
	"github.com/MrDefinite/gedis/database"
)

type typeCommandProc struct {
}

func (c typeCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {
	key := args[0]

	if key == nil {
		return types.CommonObjects.Nullbulk
	}

	obj, err := types.GetEncodeString(key)
	if err != nil {
		return types.CommonObjects.Err
	}

	return obj
}
