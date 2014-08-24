package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/edmontongo/go-zombies/server/room"
	"github.com/edmontongo/gobot/platforms/sphero"
)

var addr = flag.String("addr", ":11235", "Address to bind http server to")

// Initial implementation only supports one room
var sim = room.New()

func main() {
	flag.Parse()

	http.HandleFunc("/status", roomStatus)
	http.HandleFunc("/register", registerPlayer)
	http.HandleFunc("/collision", collidePlayer)

	fmt.Printf("Listening at %s...\n", *addr)
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
		http.Error(w, `{"error": "No name provided!"}`, http.StatusBadRequest)
		return
	}

	role := room.Human
	if req.FormValue("role") != "" {
		switch req.FormValue("role") {
		case "human":
		case "zombie":
			role = room.Zombie
		default:
			http.Error(w, fmt.Sprintf(`{"error": "Unknown role type '%s'!"`, req.FormValue("role")), http.StatusBadRequest)
		}
	}

	fmt.Fprintf(w, `{"playerId": %d}`, sim.AddPlayer(name, role, net.ParseIP(req.RemoteAddr)))
}

func collidePlayer(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := req.FormValue("id")
	if name == "" {
		http.Error(w, `{"error": "No id provided!"}`, http.StatusBadRequest)
		return
	}

	var c sphero.Collision
	data := req.FormValue("data")
	if data != "" {
		err := unwrap(data, &c)
		if err != nil {
			http.Error(w, `{"error": "Bad data!"}`, http.StatusBadRequest)
		}
	}

	id, err := room.IdFromString(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r, hit, err := sim.Collision(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, `{"role": "%s", "hit": "%s"}`, r, hit)
}

func unwrap(s string, i interface{}) error {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, i)
}
