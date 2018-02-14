package server

import (
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database/command"
	"time"
)

const (
	elNoneEvents = iota
	elAllEvents
	elFileEvents
	elTimeEvents
	elDontWait
)

const (
	feRequest = iota
	feResponse
)

type fileProc interface {
	execute(eventLoop *eventLoop, c *GedisClient)
}

type requestFileProc struct {
	req *request
}

func (rp requestFileProc) execute(eventLoop *eventLoop, c *GedisClient) {
	log.Debugf("Executing request file proc now with %d params", len(rp.req.data))
	resObj, err := command.DispatchCommand(c.DB, rp.req.data)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}
	c.enqueueResponseObj([]*basicdata.GedisObject{resObj})
}

type responseFileProc struct {
	res *response
}

func (wp responseFileProc) execute(eventLoop *eventLoop, c *GedisClient) {
	c.sendResponse()
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
		// try to get a new request
		req := client.dequeueRequest()
		if req != nil {
			log.Debug("Creating request file event now")
			req := requestFileProc{
				req: req,
			}
			event := CreateFileEvent(client, req, feRequest)
			gs.el.fileEvents = append(gs.el.fileEvents, event)
		}
		res := client.dequeueResponse()
		if res != nil {
			log.Debug("Creating response file event now")
			res := responseFileProc{
				res: res,
			}
			event := CreateFileEvent(client, res, feResponse)
			gs.el.fileEvents = append(gs.el.fileEvents, event)
		}
	}
}

func processFileEvents(gs *GedisServer) {
	for i := len(gs.el.fileEvents) - 1; i >= 0; i-- {
		fe := gs.el.fileEvents[i]
		var proc fileProc
		if fe.eventType == feRequest {
			proc = fe.requestFileProc
		} else if fe.eventType == feResponse {
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
	if flags == elAllEvents {
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

	if mask == feRequest {
		fe.requestFileProc = proc
	}
	if mask == feResponse {
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
		processEvents(gs, elAllEvents)
	}

}
