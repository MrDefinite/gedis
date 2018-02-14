package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
)

type cmdCommandProc struct {
}

func (c cmdCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {

	return nil
}
