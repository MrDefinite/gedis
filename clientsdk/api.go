package clientsdk

import (
	"github.com/sirupsen/logrus"
	"net"
	"errors"
	"sync"
	"strconv"
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



