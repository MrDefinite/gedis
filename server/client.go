package server

import (
	"net"
	"github.com/google/uuid"
	"io"
	"fmt"
	"github.com/MrDefinite/gedis/database"
	"strings"
)

type GedisClient struct {
	conn       net.Conn
	db         *database.GedisDB
	clientName string
	clientId   uuid.UUID
	cmdArgs    []string
}

func createClient(conn net.Conn) *GedisClient {
	log.Debug("Creating new client for incoming connection...")

	clientId := uuid.New()
	log.Debugf("Allocate ID '%s' for the new client, which is from '%s'",
		clientId.String(), conn.RemoteAddr().String())

	gs := GedisClient{conn: conn, clientId: clientId, cmdArgs: nil}

	return &gs
}

func (gs *GedisClient) receiveCmd() {
	log.Debugf("Listening to cmd for client '%s'", gs.clientId.String())

	var (
		buff = make([]byte, 1024)
	)
	for {
		n, err := gs.conn.Read(buff)
		data := string(buff[:n])

		switch err {
		case io.EOF:
			log.Debugf("Connection closed for client '%s'", gs.clientId.String())
			return
		case nil:
			log.Debugf("Received cmd is: %s", data)
			gs.cmdArgs = parseCmd(data)
		default:
			fmt.Errorf("receive data failed: %s", err)
			return
		}
	}
}

func (gs *GedisClient) sendResponse() {
	gs.conn.Write([]byte("test"))
}

func parseCmd(cmd string) []string {
	cmd = strings.TrimSpace(cmd)
	args := strings.Split(cmd, " ")
	return args
}
