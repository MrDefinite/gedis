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

var (
	ErrAnotherRequestProcessing = errors.New("there is another request being processed")
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
	writer *resp.Writer
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
	gc.parser = resp.CreateNewParser(conn)
	gc.writer = resp.CreateNewWriter(conn)

	return nil
}

func (gc *Gclient) CloseConnection() error {
	if gc.conn == nil {
		return errors.New("there is no connection there")
	}

	gc.conn.Close()

	return nil
}

func (gc *Gclient) heartbeat() {

}

func (gc *Gclient) SetConnectionTimeout() {
	gc.conn.SetReadDeadline(time.Now().Add(time.Duration(gc.requestTimeout) * time.Second))
}

// Return the output string, and the error if there is
func (gc *Gclient) ProcessCmdString(cmds []string) (string, error) {
	if gc.isCommunicatingToServer {
		return "", ErrAnotherRequestProcessing
	}

	gc.commLock.Lock()
	gc.isCommunicatingToServer = true
	defer gc.commLock.Unlock()

	// pos 0 is cmd, and pos 1,2,3... is args
	arrayLen := len(cmds)
	gc.writer.AppendArrayLength(arrayLen)

	for _, d := range cmds {
		gc.writer.AppendBulkString(d)
	}

	gc.SetConnectionTimeout()
	gc.writer.Write()

	// Now let's wait for the response
	pr, err := gc.parser.Parse()
	if err != nil {
		return "", err
	}

	res, err := gc.parser.FormatResultAsString(pr)
	if err != nil {
		return "", err
	}

	gc.isCommunicatingToServer = false
	return res, nil
}
