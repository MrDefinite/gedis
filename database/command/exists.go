package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/protocol/resp"
)

type existsCommandProc struct {
}

func (c existsCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]

	if key == nil {
		return nil
	}

	return resp.CommonObjects.Ok
}
