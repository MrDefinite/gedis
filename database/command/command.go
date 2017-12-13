package command

import (
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/database/types"
)

type cmdCommandProc struct {
}

func (c cmdCommandProc) execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject {

	return nil
}
