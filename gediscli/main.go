package main

import (
	"log"
	"flag"
	"net"
	"strconv"
	"bufio"
	"os"
)


const (
	DEFAULT_CONN_TYPE = "tcp"
	DEFAULT_CONN_PORT = 9019
	DEFAULT_CONN_HOST = "127.0.0.1"
)


func main() {
	log.Println("Initializing gedis cli service...")

	server := flag.String("server", DEFAULT_CONN_HOST, "gedis server address")
	port := flag.Int("port", DEFAULT_CONN_PORT, "gedis server port")
	gedisServer := *server + ":" + strconv.Itoa(*port)

	log.Printf("Connecting to gedis server '%s'", gedisServer)

	conn, err := net.Dial(DEFAULT_CONN_TYPE, gedisServer)
	if err != nil {
		log.Fatalln("Failed to connect to gedis server, the error is: " + err.Error())
	}

	log.Println("Initializing console reader...")
	reader := bufio.NewReader(os.Stdin)
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalln("Failed to read command line from console, the error is: " + err.Error())
			break
		}
		line := string(bytes)
		log.Println("Incoming command is: " + line)
	}

	conn.Close()
}

