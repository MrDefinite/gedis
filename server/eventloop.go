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

func processFileEvents() {

}

func processTimeEvents() {

}

func processEvents(eventLoop *eventLoop) {

	processFileEvents()
	processTimeEvents()

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
		processEvents(gs.el)
	}

}
