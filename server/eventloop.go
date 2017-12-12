package server

import (
	"time"
	"github.com/MrDefinite/gedis/database/command"
	"github.com/MrDefinite/gedis/database/types"
)

const (
	EL_NONE_EVENTS = 0
	EL_ALL_EVENTS
	EL_FILE_EVENTS
	EL_TIME_EVENTS
	EL_DONT_WAIT
)

const (
	FE_REQUEST  = 0
	FE_RESPONSE
)

type fileProc interface {
	execute(eventLoop *eventLoop, c *GedisClient)
}

type requestFileProc struct {
}

func (rp requestFileProc) execute(eventLoop *eventLoop, c *GedisClient) {
	log.Debugf("Executing request file proc now with %d params", len(c.CmdArgs))
	resObj, err := command.DispatchCommand(c.DB, c.CmdArgs)
	if err != nil {
		log.Errorf("%s", err.Error())
		c.CmdArgs = nil
		return
	}
	c.CmdArgs = nil
	c.Response = resObj
}

type responseFileProc struct {
}

func (wp responseFileProc) execute(eventLoop *eventLoop, c *GedisClient) {
	response := c.Response

	c.sendResponse(types.GetStringValueFromObject(response))

	c.Response = nil
}

type fileEvent struct {
	eventType        uint8
	requestFileProc  fileProc
	responseFileProc fileProc
	client           *GedisClient
}

type timeEvent struct {
}

type eventLoop struct {
	timeEventNextId int64
	lastTime        time.Time

	fileEvents []*fileEvent
	timeEvents []*timeEvent

	stop bool
}

func CreateEventLoop(gs *GedisServer) {
	log.Info("Creating main event loop now")
	gs.el = &eventLoop{}
}

func checkAndCreateEvents(gs *GedisServer) {
	for _, client := range gs.clients {
		if client.CmdArgs != nil {
			req := requestFileProc{}
			log.Debug("Creating request file event now")
			event := CreateFileEvent(client, req, FE_REQUEST)
			gs.el.fileEvents = append(gs.el.fileEvents, event)
		}
		if client.Response != nil {
			res := responseFileProc{}
			log.Debug("Creating response file event now")
			event := CreateFileEvent(client, res, FE_RESPONSE)
			gs.el.fileEvents = append(gs.el.fileEvents, event)
		}
	}
}

func processFileEvents(gs *GedisServer) {
	for i := len(gs.el.fileEvents) - 1; i >= 0; i-- {
		fe := gs.el.fileEvents[i]
		var proc fileProc
		if fe.eventType == FE_REQUEST {
			proc = fe.requestFileProc
		} else if fe.eventType == FE_RESPONSE {
			proc = fe.responseFileProc
		}
		if proc != nil {
			proc.execute(gs.el, fe.client)
			gs.el.fileEvents = append(gs.el.fileEvents[:i],
				gs.el.fileEvents[i+1:]...)
		}
	}
}

func processTimeEvents(gs *GedisServer) {

}

func processEvents(gs *GedisServer, flags uint8) {
	if flags == EL_ALL_EVENTS {
		processFileEvents(gs)
		processTimeEvents(gs)
	}
}

func CreateFileEvent(client *GedisClient, proc fileProc, mask uint8) *fileEvent {
	fe := fileEvent{
		eventType:        mask,
		requestFileProc:  nil,
		responseFileProc: nil,
		client:           client,
	}

	if mask == FE_REQUEST {
		fe.requestFileProc = proc
	}
	if mask == FE_RESPONSE {
		fe.responseFileProc = proc
	}

	return &fe
}

func createTimeEvent() {

}

func StopLoop(eventLoop *eventLoop) {
	eventLoop.stop = true
}

func DeleteEventLoop(eventLoop *eventLoop) {

}

func MainLoop(gs *GedisServer) {
	gs.el.stop = false

	for gs.el.stop != true {
		// Check and create file events
		checkAndCreateEvents(gs)

		// Run event loop
		processEvents(gs, EL_ALL_EVENTS)
	}

}
