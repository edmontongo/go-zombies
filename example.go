package main

import (
	"flag"
	"log"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

/*
	Mac:
	Pair with a Bluetooth device, then specify a port in the form:
 		/dev/tty.Sphero-???-RN-SPP
 	where ??? are the colours when pairing (eg. ROG).

 	Windows:
 	Pair a Bluetooth device and look up the COM port, use the form:
 		COM2

	Linux:
	Use rfcomm ... and then specify the port in the form:
		/dev/rfcomm0

	Sphero documentation:
	http://gobot.io/documentation/platforms/sphero/#HowToConnect
*/
var device = flag.String("device", "/dev/tty.Sphero-WRW-RN-SPP", "Device for the Sphero.")

var server = flag.String("server", "http://localhost:11235", "Server address to connect to.")

// true to start as a zombie
var zombie = flag.Bool("zombie", false, "Runs the example as a zombie.")

func zombieTicker(zombie game.Robot) {
	zombie.SetReferenceHeading(0)
	c := time.Tick(1 * time.Second)
	go func() {
		heading := 0
		for {
			select {
			case <-c:
				zombie.Walk(100, heading)
				heading += 6
			case event, ok := <-zombie.Events:
				if !ok {
					return
				}
				log.Printf("Event %+v\n.", event)
				heading += 180
				zombie.Walk(100, heading)
			}
		}
	}()
}

func main() {
	flag.Parse()

	game.RegisterZombie(zombieTicker)
	game.RegisterHuman(zombieTicker)

	err := game.Start("bob", *zombie, *device, *server)
	if err != nil {
		log.Fatal(err)
	}
}
