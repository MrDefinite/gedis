package server

import (
	"net"
	"github.com/google/uuid"
	"io"
	"fmt"
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/database/types"
	//"github.com/MrDefinite/gedis/common/protocol/resp"
)

type GedisClient struct {
	conn       net.Conn
	DB         *database.GedisDB
	clientName string
	clientId   uuid.UUID
	CmdArgs    []*types.GedisObject
	Response   *types.GedisObject
	running    bool
}

func createClient(conn net.Conn, db *database.GedisDB) *GedisClient {
	log.Debug("Creating new client for incoming connection...")

	clientId := uuid.New()
	log.Debugf("Allocate ID '%s' for the new client, which is from '%s'",
		clientId.String(), conn.RemoteAddr().String())

	c := GedisClient{
		conn:     conn,
		clientId: clientId,
		CmdArgs:  nil,
		Response: nil,
		running:  true,
		DB:       db,
	}

	// If conn is nil, then it is fake client
	if conn != nil {
		go c.handleCmd()
	}

	return &c
}

func tearDownClient(c *GedisClient) {
	c.running = false
	conn := c.conn
	conn.Close()
	c = nil
}

func (c *GedisClient) handleCmd() {
	log.Debugf("Listening to cmd for client '%s'", c.clientId.String())

	var (
		buff = make([]byte, 1024)
	)
	for c.running == true {
		n, err := c.conn.Read(buff)
		data := string(buff[:n])

		switch err {
		case io.EOF:
			log.Debugf("Connection closed for client '%s'", c.clientId.String())
			return
		case nil:
			if c.CmdArgs != nil {
				log.Warnf("Already commands being processed!")
				break
			}

			log.Debugf("Received cmd is: %q\n", data)
			c.CmdArgs, err = parseCmd(data)
			if err != nil {
				c.Response = types.CommonObjects.Syntaxerr
			}
		default:
			fmt.Errorf("receive data failed: %s", err)
			return
		}
	}
}

func (c *GedisClient) sendResponse(response string) {
	c.conn.Write([]byte(response))
}

func parseCmd(cmd string) ([]*types.GedisObject, error) {
	//requestCmd, err := resp.ParseRequest(cmd)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var objArgs []*types.GedisObject
	//objArgs = append(objArgs, types.CreateStringObject(requestCmd.Name))
	//
	//for _, arg := range requestCmd.Params {
	//	objArg := types.CreateStringObject(arg)
	//	objArgs = append(objArgs, objArg)
	//}

	return nil, nil
}
