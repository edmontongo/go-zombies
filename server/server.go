package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/edmontongo/go-zombies/server/room"
)

var addr = flag.String("addr", ":11235", "Address to bind http server to")

// Initial implementation only supports one room
var sim = room.New()

func main() {
	flag.Parse()

	http.HandleFunc("/status", roomStatus)
	http.HandleFunc("/register", registerPlayer)
	http.HandleFunc("/collision", collidePlayer)
	http.HandleFunc("/deregister", deregisterPlayer)

	fmt.Printf("Listening at %s...\n", *addr)
	panic(http.ListenAndServe(*addr, nil))
}

var reportTemplate = `<html><head><title>Zombie Simulartor Status</title></head>
<body>
<div id="Stats">Humans: {{.Humans}} Zombies: {{.Zombies}}<div>
<br>
<div id="players">Players:<br>
{{range .Players}}{{.Description}}<br>{{end}}
</div>
<br>
<div id="recent">Recent Collisions:<br>
{{range .Recent}}
{{with .}}{{.}}<br>{{end}}
{{end}}
</div>
</body></html>
`
var statusTemplate = template.Must(template.New("reportTemplate").Parse(reportTemplate))

// roomStatus provides very basic status information
func roomStatus(w http.ResponseWriter, req *http.Request) {
	err := statusTemplate.Execute(w, &sim)
	if err != nil {
		log.Println(err.Error())
	}
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

	id := sim.AddPlayer(name, role, net.ParseIP(req.RemoteAddr))

	log.Printf("Player '%s' from %v given id %d as %s", name, req.RemoteAddr, id, role)
	fmt.Fprintf(w, `{"playerId": %d}`, id)
}

func deregisterPlayer(w http.ResponseWriter, req *http.Request) {
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
	id, err := room.IdFromString(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sim.RemovePlayer(id)
}

func collidePlayer(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := req.FormValue("id")
	if id == "" {
		http.Error(w, `{"error": "No id provided!"}`, http.StatusBadRequest)
		return
	}

	var c room.Collision
	data := req.FormValue("data")
	if data != "" {
		err := unwrap(data, &c.Collision)
		if err != nil {
			http.Error(w, `{"error": "Bad data!"}`, http.StatusBadRequest)
		}
		log.Printf("Collision from %s: %v", id, c.Collision)
	}

	c.Id, err = room.IdFromString(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ServerTime = time.Now()

	r, hit, err := sim.Collision(c)
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
