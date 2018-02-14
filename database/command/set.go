package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
)

type setCommandProc struct {
}

func (c setCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]
	value := args[1]

	value = basicdata.TryObjectEncoding(value)

	db.Dict[*key] = value

	return basicdata.CommonObjects.Ok
}
