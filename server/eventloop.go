package server

import (
	"time"
)

type FileEvent struct {
}

type TimeEvent struct {
}

type FiredEvent struct {
}

type eventLoop struct {
	maxConnection   int
	setSize         int
	timeEventNextId int64
	lastTime        time.Time

	fileEvents *FileEvent
	fired      *FiredEvent
	timeEvents *TimeEvent

	stop bool
}

func CreateEventLoop(gs *GedisServer) {
	log.Info("Creating main event loop now")
	gs.el = &eventLoop{}
}

func processFileEvents(gs *GedisServer) {
	for _, client := range gs.clients {
		if client.cmdArgs != nil {

			client.sendResponse()
			client.cmdArgs = nil
		}
	}
}

func processTimeEvents(gs *GedisServer) {

}

func processEvents(gs *GedisServer) {
	processFileEvents(gs)
	processTimeEvents(gs)
}

func CreateFileEvent() {

}

func CreateTimeEvent() {

}

func StopLoop(eventLoop *eventLoop) {
	eventLoop.stop = true
}

func DeleteEventLoop(eventLoop *eventLoop) {

}

func MainLoop(gs *GedisServer) {
	gs.el.stop = false

	for gs.el.stop != true {
		// Run event loop
		processEvents(gs)
	}

}
