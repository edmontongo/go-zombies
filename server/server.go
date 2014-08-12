package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/edmontongo/go-zombies/server/room"
)

var addr = flag.String("addr", ":11235", "Address to bind http server to")

// Initial implementation only supports one room
var sim = room.New()

func main() {
	flag.Parse()

	http.HandleFunc("/status", roomStatus)
	http.HandleFunc("/register", registerPlayer)

	panic(http.ListenAndServe(*addr, nil))
}

// roomStatus provides very basic status information
func roomStatus(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Zombie Simulartor Status</title></head><body>\n")
	fmt.Fprintf(w, "Humans: %d<br>Zombies: %d<br>\n", sim.Humans(), sim.Zombies())
	fmt.Fprintf(w, "</body></html>\n")
}

func registerPlayer(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := req.FormValue("name")
	if name == "" {
		http.Error(w, "No name provided!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, `{"playerId": %d}`, sim.AddPlayer(name, room.Human, net.ParseIP(req.RemoteAddr)))
}
