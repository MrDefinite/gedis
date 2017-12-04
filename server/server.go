package server

import (
	"log"
	"github.com/MrDefinite/gedis/database"
	"net"
	"os"
	"strconv"
)

const (
	DEFAULT_CONN_TYPE = "tcp"
	DEFAULT_CONN_PORT = 9019
	DEFAULT_CONN_HOST = "0.0.0.0"
)

type GedisServer struct {
	address        string
	port           int
	connectionType string

	listener net.Listener

	el *eventLoop

	db *database.GedisDB

	clients []*GedisClient
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
			log.Println("Error accepting: ", err.Error())
			break
		}
		// Create client for incoming connection
		client := GedisClient{conn:conn}
		gs.clients = append(gs.clients, &client)
	}
}

func initDB() {
	log.Println("Initializing gedis databases now")
}

func InitServer(gs *GedisServer) {
	log.Println("Initializing gedis gs now")

	CreateEventLoop(gs)

	initDB()

	listenToPort(gs)
	go handleConnections(gs)


}

func InitServerConfig(gs *GedisServer) {
	log.Println("Initializing gedis server configuration now")

	gs.address = DEFAULT_CONN_HOST
	gs.port = DEFAULT_CONN_PORT
	gs.connectionType = DEFAULT_CONN_TYPE
}

func TearDownServer(gs *GedisServer) {
	gs.listener.Close()

	DeleteEventLoop(gs.el)
}
