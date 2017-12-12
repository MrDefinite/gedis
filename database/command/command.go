package command

import (
	"fmt"
	"github.com/MrDefinite/gedis/database/types"
	"github.com/pkg/errors"
	"github.com/MrDefinite/gedis/database"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

var commandDict = map[string]*Command{
	"get":  initCommand("get", getCommandProc{}),
	"set":  initCommand("set", setCommandProc{}),
	"type": initCommand("type", typeCommandProc{}),
}

type Command struct {
	name         string
	proc         commandProc
	argc         int
	firstKey     int
	lastKey      int
	keyStep      int
	microseconds int64
	calls        int64
}

type commandProc interface {
	execute(db *database.GedisDB, args []*types.GedisObject) *types.GedisObject
}

func initCommand(name string, proc commandProc) *Command {
	cmd := Command{
		name:         name,
		proc:         proc,
		argc:         0,
		firstKey:     0,
		lastKey:      0,
		keyStep:      0,
		microseconds: 0,
		calls:        0,
	}
	return &cmd
}

func DispatchCommand(db *database.GedisDB, cmdArgs []*types.GedisObject) (*types.GedisObject, error) {
	log.Debugf("Dispatching the command...")
	commandName := types.GetStringValueFromObject(cmdArgs[0])
	currentCommand := commandDict[commandName]
	if currentCommand == nil {
		fmt.Errorf("cannot find command named as: %s", commandName)
		return nil, errors.New("cannot find command named as: " + commandName)
	}

	args := cmdArgs[1:]
	proc := currentCommand.proc
	return proc.execute(db, args), nil
}
