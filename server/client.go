package server

import (
	"net"
	"github.com/google/uuid"
)

type GedisClient struct {
	conn net.Conn
	dbId int
	clientName string
	clientId uuid.UUID
}

func CreateClient(conn net.Conn) *GedisClient {
	log.Debug("Creating new client for incoming connection...")

	clientId := uuid.New()
	log.Debugf("Allocate ID '%s' for the new client", clientId.String())

	return &GedisClient{conn: conn, clientId: clientId}
}

