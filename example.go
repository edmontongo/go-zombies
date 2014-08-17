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
var port = flag.String("port", "/dev/tty.Sphero-OBY-RN-SPP", "Port to the Sphero")

// true to start as a zombie
var zombie = flag.Bool("zombie", false, "Runs the example as a zombie")

func zombieTicker(zombie game.Robot) {
	c := time.Tick(1 * time.Second)
	go func() {
		heading := 0
		for {
			select {
			case <-c:
				zombie.Walk(10, heading)
				heading += 6
			case event, ok := <-zombie.Events:
				if !ok {
					return
				}
				log.Printf("Event %v\n.", event)
				heading += 180
				zombie.Walk(10, heading)
			}
		}
	}()
}

func main() {
	flag.Parse()

	game.RegisterZombie(zombieTicker)
	// game.RegisterHuman(me)

	err := game.Start("bob", *zombie, *port)
	if err != nil {
		log.Fatal(err)
	}
}
