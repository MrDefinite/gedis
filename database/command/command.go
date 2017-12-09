package command

import (
	"github.com/MrDefinite/gedis/server"
	"fmt"
)

var commandDict = map[string]*Command {
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
	execute(client *server.GedisClient)
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

func dispatchCommand(cmd string, client *server.GedisClient) {
	currentCommand := commandDict[cmd]
	if currentCommand == nil {
		fmt.Errorf("cannot find command named as: %s", cmd)
		return
	}

	proc := currentCommand.proc
	proc.execute(client)
}

