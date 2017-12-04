package main

import "log"
import (
	"github.com/MrDefinite/gedis/server"
)

func main() {
	log.Println("Gedis begins to run!")

	gs := &server.GedisServer{}

	server.InitServerConfig(gs)
	server.InitServer(gs)


	server.MainLoop(gs)

	server.TearDownServer(gs)


	log.Println("Gedis Stopped!")
}
