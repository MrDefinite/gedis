package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
)

type getCommandProc struct {
}

func (c getCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]
	value := db.Dict[*key]

	if value == nil {
		return basicdata.CommonObjects.NullBulk
	}
	return value.(*basicdata.GedisObject)
}
