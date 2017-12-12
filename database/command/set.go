package command

import (
	"github.com/MrDefinite/gedis/database/types"
	"github.com/MrDefinite/gedis/database"
)

type setCommandProc struct {
}

func (c setCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {
	key := args[0]
	value := args[1]

	value = types.TryObjectEncoding(value)

	db.Dict[*key] = value

	return types.CommonObjects.Ok
}
