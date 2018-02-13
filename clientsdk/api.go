package clientsdk

import (
	"errors"
	"github.com/MrDefinite/gedis/protocol/resp"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	log = logrus.New()
)

const (
	defaultConnType = "tcp"
	defaultLogLevel = 3
	// Default write to console
	defaultLogOutput = ""
	// Disable idle timeout by default
	defaultIdleTimeout = -1
	// 20 seconds
	defaultRequestTimeout = 20
)

type Gclient struct {
	connectionType          string
	logLevel                logrus.Level
	conn                    net.Conn
	idleTimeout             int64
	Server                  string
	requestTimeout          int32
	isCommunicatingToServer bool
	commLock                sync.Mutex

	parser *resp.Parser
}

func CreateNewInstance() *Gclient {
	gc := Gclient{}
	gc.connectionType = defaultConnType

	gc.idleTimeout = defaultIdleTimeout
	gc.requestTimeout = defaultRequestTimeout

	gc.isCommunicatingToServer = false
	gc.commLock = sync.Mutex{}

	gc.SetLogLevel(defaultLogLevel)
	gc.SetLogOutput(defaultLogOutput)

	return &gc
}

func (gc *Gclient) ConnectToServer(host string, port int) error {
	if host == "" || port <= 0 {
		return errors.New("illegal address or port for gedis server")
	}

	gc.Server = host + ":" + strconv.Itoa(port)
	conn, err := net.Dial(defaultConnType, gc.Server)
	if err != nil {
		return err
	}

	gc.conn = conn

	return nil
}

func (gc *Gclient) CloseConnection() error {
	if gc.conn == nil {
		return errors.New("there is no connection there")
	}

	gc.conn.Close()

	return nil
}

func (gc *Gclient) ParseAndProcessCmd(cmd string) (string, error) {

	return "", nil
}

func (gc *Gclient) heartbeat() {

}

func (gc *Gclient) sendRequestAndGetResponse(encodedRequest []byte) ([]byte, error) {
	if encodedRequest == nil {
		return nil, errors.New("cannot send empty request to server")
	}

	if gc.isCommunicatingToServer {
		return nil, errors.New("there is another request being processed")
	}

	gc.commLock.Lock()

	// Send it to server now
	gc.conn.Write(encodedRequest)

	// Init buffer
	buff := make([]byte, resp.MaxDataSizeReadPerTime)

	gc.conn.SetReadDeadline(time.Now().Add(time.Duration(gc.requestTimeout) * time.Second))
	// Wait for response from server

	var n int
	var err error
	//var parser resp.
	for n, err = gc.conn.Read(buff); ; {
		if err != nil {
			gc.commLock.Unlock()
			return nil, err
		}

	}

	gc.commLock.Unlock()

	return buff[:n], nil
}

// Return the output string, and the error if there is
func (gc *Gclient) ProcessCmdString(cmds []string) (string, error) {

}
