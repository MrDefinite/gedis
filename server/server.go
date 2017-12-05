package server

import (
	"github.com/MrDefinite/gedis/database"
	"net"
	"os"
	"strconv"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_CONN_TYPE = "tcp"
	DEFAULT_CONN_PORT = 9019
	DEFAULT_CONN_HOST = "0.0.0.0"
)

var (
	log = logrus.New()
)

type GedisServer struct {
	address        string
	port           int
	connectionType string

	listener net.Listener

	el *eventLoop

	db *database.GedisDB

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
		client := createClient(conn)
		gs.clients = append(gs.clients, client)
		go client.receiveCmd()
	}
}

func initDB() {
	log.Info("Initializing gedis databases now")
}

func initLogger(gs *GedisServer) {
	log.SetLevel(gs.logLevel)
}

func InitServer(gs *GedisServer) {
	log.Info("Initializing gedis gs now")

	initLogger(gs)

	CreateEventLoop(gs)

	initDB()

	listenToPort(gs)
	go handleConnections(gs)


}

func InitServerConfig(gs *GedisServer) {
	log.Info("Initializing gedis server configuration now")

	gs.address = DEFAULT_CONN_HOST
	gs.port = DEFAULT_CONN_PORT
	gs.connectionType = DEFAULT_CONN_TYPE
	gs.logLevel = logrus.DebugLevel
}

func TearDownServer(gs *GedisServer) {
	gs.listener.Close()

	DeleteEventLoop(gs.el)
}
