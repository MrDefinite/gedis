package command

import (
	"github.com/MrDefinite/gedis/database/types"
	"github.com/MrDefinite/gedis/database"
)

type typeCommandProc struct {
}

func (c typeCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {

	return nil
}
