package server

import (
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/database/types"
	"net"
	"os"
	"strconv"
	"github.com/sirupsen/logrus"
)

const (
	defaultConnType   = "tcp"
	defaultConnPort   = 9019
	defaultConnHost   = "0.0.0.0"
	defaultMaxClients = 10000
)

var (
	log = logrus.New()
)

type GedisServer struct {
	address        string
	port           int
	connectionType string

	maxClients int

	listener net.Listener

	el *eventLoop

	db []*database.GedisDB

	clients []*GedisClient

	logLevel logrus.Level
}

func listenToPort(gs *GedisServer) {
	address := gs.address
	port := strconv.Itoa(gs.port)

	log.Println("Bind listener to port: " + port)
	l, err := net.Listen(gs.connectionType, address+":"+port)
	if err != nil {
		log.Fatalln("Error listening: ", err.Error())
		os.Exit(1)
	}

	gs.listener = l
}

func handleConnections(gs *GedisServer) {
	for {
		// Listen for an incoming connection.
		conn, err := gs.listener.Accept()
		if err != nil {
			log.Warn("Error accepting: ", err.Error())
			break
		}
		// Create client for incoming connection
		client := createClient(conn, gs.db[0])
		if len(gs.clients) >= gs.maxClients {
			log.Warn("Max client number exceeds")
			client.sendResponse("Max client number exceeds, please connect later")
			tearDownClient(client)
			continue
		}

		log.Debug("Enqueue new client instance")
		gs.clients = append(gs.clients, client)
	}
}

func initDB(gs *GedisServer) {
	log.Info("Initializing gedis databases now")

	gs.db = append(gs.db, database.InitDBFromZero())
}

func initLogger(gs *GedisServer) {
	log.SetLevel(gs.logLevel)
}

func InitServer(gs *GedisServer) {
	log.Info("Initializing gedis gs now")

	initLogger(gs)

	types.InitCommonObjects()

	CreateEventLoop(gs)

	initDB(gs)

	listenToPort(gs)
	go handleConnections(gs)

}

func InitServerConfig(gs *GedisServer) {
	log.Info("Initializing gedis server configuration now")

	gs.address = defaultConnHost
	gs.port = defaultConnPort
	gs.connectionType = defaultConnType
	gs.maxClients = defaultMaxClients
	gs.logLevel = logrus.DebugLevel
}

func TearDownServer(gs *GedisServer) {
	gs.listener.Close()

	DeleteEventLoop(gs.el)
}
