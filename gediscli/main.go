package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/MrDefinite/gedis/clientsdk"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	defaultConnPort = 9019
	defaultConnHost = "127.0.0.1"
	defaultLogLevel = logrus.DebugLevel
)

var (
	log = logrus.New()
)

func main() {
	log.Level = defaultLogLevel
	log.Out = os.Stdout

	log.Info("Initializing gedis cli service...")
	log.Info("Creating gedis client service...")

	gClient := clientsdk.CreateNewInstance()

	host := flag.String("host", defaultConnHost, "gedis host address")
	port := flag.Int("port", defaultConnPort, "gedis host port")

	err := gClient.ConnectToServer(*host, *port)
	if err != nil {
		log.Fatalln("Failed to connect to gedis host, the error is: " + err.Error())
	}

	log.Info("Connected to server, initializing console reader now...")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s > ", gClient.Server)
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Error("Failed to read command line from console, the error is: " + err.Error())
			break
		}
		line := string(bytes)

		if simpleCmdCheck(line) == false {
			log.Warnf("Command %s is not a valid!", line)
			continue
		}

		cmds := formatCmd(line)
		outputStr, err := gClient.ProcessCmdString(cmds)
		if err != nil {
			log.Error(err.Error())
			break
		}

		fmt.Println("\"" + outputStr + "\"")
	}

	gClient.CloseConnection()
}

func simpleCmdCheck(cmd string) bool {
	if len(cmd) == 0 {
		return false
	}

	return true
}

func formatCmd(cmd string) []string {
	strs := strings.Split(cmd, " ")
	var cmds []string
	for _, s := range strs {
		s = strings.TrimSpace(s)
		cmds = append(cmds, s)
	}

	return cmds
}
