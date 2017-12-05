package main

import (
	"flag"
	"net"
	"strconv"
	"bufio"
	"os"
	"github.com/sirupsen/logrus"
)


const (
	DEFAULT_CONN_TYPE = "tcp"
	DEFAULT_CONN_PORT = 9019
	DEFAULT_CONN_HOST = "127.0.0.1"
	DEFAULT_LOG_LEVEL = logrus.DebugLevel
)

var (
	log = logrus.New()
)

func main() {
	log.Level = DEFAULT_LOG_LEVEL
	log.Out = os.Stdout

	log.Info("Initializing gedis cli service...")

	server := flag.String("server", DEFAULT_CONN_HOST, "gedis server address")
	port := flag.Int("port", DEFAULT_CONN_PORT, "gedis server port")
	gedisServer := *server + ":" + strconv.Itoa(*port)

	log.Infof("Connecting to gedis server '%s'", gedisServer)

	conn, err := net.Dial(DEFAULT_CONN_TYPE, gedisServer)
	if err != nil {
		log.Fatalln("Failed to connect to gedis server, the error is: " + err.Error())
	}

	log.Info("Connected to server, initializing console reader now...")
	reader := bufio.NewReader(os.Stdin)
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Error("Failed to read command line from console, the error is: " + err.Error())
			break
		}
		line := string(bytes)
		log.Debug("Incoming command is: " + line)

		if simpleCmdCheck(line) == false {
			log.Warnf("Command %s is not a valid!", line)
		}

		// Send it to server now
		conn.Write([]byte(line))

		// Wait for response from server
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			log.Error("Failed to read response from server, the error is: " + err.Error())
			continue
		}
		log.Debugf("Receive: %s", buff[:n])
	}

	conn.Close()
}

func simpleCmdCheck(cmd string) bool {
	if len(cmd) == 0 {
		return false
	}

	return true
}
