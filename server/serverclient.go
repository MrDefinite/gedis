package server

import (
	"fmt"
	"github.com/MrDefinite/gedis/basicdata"
	"github.com/MrDefinite/gedis/database"
	"github.com/MrDefinite/gedis/protocol/resp"
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
)

type request struct {
	data []*basicdata.GedisObject
}

type response struct {
	data []*basicdata.GedisObject
}

type GedisClient struct {
	conn       net.Conn
	DB         *database.GedisDB
	clientName string
	clientId   uuid.UUID

	reqLock sync.Mutex
	resLock sync.Mutex

	reqs  []*request
	resps []*response

	running bool

	parser *resp.Parser
	writer *resp.Writer
}

func createClient(conn net.Conn, db *database.GedisDB) *GedisClient {
	log.Debug("Creating new client for incoming connection...")

	clientId := uuid.New()
	log.Debugf("Allocate ID '%s' for the new client, which is from '%s'",
		clientId.String(), conn.RemoteAddr().String())

	c := GedisClient{
		conn:     conn,
		clientId: clientId,
		running:  true,
		reqs:     make([]*request, 0),
		resps:    make([]*response, 0),
		DB:       db,
	}

	// If conn is nil, then it is fake client
	if conn != nil {
		c.parser = resp.CreateNewParser(conn)
		c.writer = resp.CreateNewWriter(conn)
	}

	return &c
}

func tearDownClient(c *GedisClient) {
	c.running = false
	conn := c.conn
	conn.Close()
	c = nil
}

func (c *GedisClient) handleRequest() {
	log.Debugf("Listening to cmd for client '%s'", c.clientId.String())

	for c.running == true {
		pr, err := c.parser.Parse()
		switch err {
		case io.EOF:
			log.Debugf("Connection closed for client '%s'", c.clientId.String())
			return
		case nil:
			cmdArgs, err := c.parser.FormatCmdResultAsGedisObject(pr)
			if err != nil {
				c.enqueueResponseObj([]*basicdata.GedisObject{basicdata.CommonObjects.SynTaxErr})
			} else {
				// process the command now
				c.enqueueRequestObj(cmdArgs)
			}
		default:
			fmt.Errorf("receive data failed: %s", err)
			return
		}
	}
}

func (c *GedisClient) enqueueResponseObj(objs []*basicdata.GedisObject) {
	c.enqueueResponse(&response{
		data: objs,
	})
}

func (c *GedisClient) enqueueResponse(res *response) {
	c.resLock.Lock()
	c.resps = append(c.resps, res)
	c.resLock.Unlock()
}

func (c *GedisClient) dequeueResponse() *response {
	c.resLock.Lock()
	defer c.resLock.Unlock()

	if len(c.resps) == 0 {
		return nil
	}

	res := c.resps[0]
	c.resps = c.resps[1:]
	return res
}

func (c *GedisClient) enqueueRequestObj(obj []*basicdata.GedisObject) {
	c.enqueueRequest(&request{
		data: obj,
	})
}

func (c *GedisClient) enqueueRequest(req *request) {
	c.reqLock.Lock()
	c.reqs = append(c.reqs, req)
	c.reqLock.Unlock()
}

func (c *GedisClient) dequeueRequest() *request {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	if len(c.reqs) == 0 {
		return nil
	}

	req := c.reqs[0]
	c.reqs = c.reqs[1:]
	return req
}

func (c *GedisClient) sendResponse(res *response) {
	c.sendResponseAsString(res)
}

func (c *GedisClient) sendResponseAsString(res *response) {
	d := res.data
	if len(d) == 0 {
		return
	}

	for _, obj := range d {
		t := basicdata.GetType(obj)
		switch t {
		case basicdata.GedisString:
			s, err := basicdata.GetStringValueFromObject(obj)
			if err != nil {
				// TODO: handle the error
				return
			}
			c.writer.AppendSimpleString(s)
			c.writer.Write()
		}
	}
}
