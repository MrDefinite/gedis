package command

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/protocol/resp"
)

type getCommandProc struct {
}

func (c getCommandProc) execute(db *database.GedisDB, args []*basicdata.GedisObject) *basicdata.GedisObject {
	key := args[0]
	value := db.Dict[*key]

	if value == nil {
		return resp.CommonObjects.NullBulk
	}

	switch dp := value.(type) {
	case *basicdata.GedisObject:
		return resp.Encoder.EncodeToResponseObj(dp)
	}

	return resp.CommonObjects.NullBulk
}
