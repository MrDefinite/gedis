package server

import "net"

type GedisClient struct {
	conn net.Conn
}

func CreateClient() *GedisClient {





	return &GedisClient{}
}

