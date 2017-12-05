package main

import (
	"github.com/MrDefinite/gedis/server"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func main() {
	log.Info("Gedis begins to run!")

	gs := &server.GedisServer{}

	server.InitServerConfig(gs)
	server.InitServer(gs)


	server.MainLoop(gs)

	server.TearDownServer(gs)


	log.Info("Gedis Stopped!")
}
