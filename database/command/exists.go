package command


import (
	"github.com/MrDefinite/gedis/database/types"
	"github.com/MrDefinite/gedis/database"
)

type existsCommandProc struct {
}

func (c existsCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {
	key := args[0]

	if key == nil {
		return nil
	}

	return types.CommonObjects.Ok
}
