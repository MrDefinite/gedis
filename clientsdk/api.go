package clientsdk

import (
	"github.com/sirupsen/logrus"
	"net"
	"errors"
	"sync"
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
	server                  string
	requestTimeout          int64
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

func (gc *Gclient) ConnectToServer(address string, port string) error {
	if address == "" || port == "" {
		return errors.New("illegal address or port for gedis server")
	}

	gc.server = address + ":" + port
	conn, err := net.Dial(defaultConnType, gc.server)
	if err != nil {
		return err
	}

	gc.conn = conn

	return nil
}



