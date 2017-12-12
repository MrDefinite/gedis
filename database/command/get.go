package command

import (
	"github.com/MrDefinite/gedis/database/types"
	"github.com/MrDefinite/gedis/database"
)

type getCommandProc struct {
}

func (c getCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {
	key := args[0]
	value := db.Dict[*key]

	if value == nil {
		return types.CommonObjects.Nullbulk
	}
	return value.(*types.GedisObject)
}
