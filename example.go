package main

import (
	"flag"
	"log"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

/*
	Mac:
	Pair with a Bluetooth device, then specify a device in the form:
 		/dev/tty.Sphero-???-RN-SPP
 	where ??? are the colours when pairing (eg. ROG).

 	Windows:
 	Pair a Bluetooth device and look up the COM port, use the form:
 		COM2

	Linux:
	Use rfcomm ... and then specify the device in the form:
		/dev/rfcomm0

	Sphero documentation:
	http://gobot.io/documentation/platforms/sphero/#HowToConnect
*/
var device = flag.String("device", "/dev/tty.Sphero-RBG-RN-SPP", "Device for the Sphero.")

// game server address
var server = flag.String("server", "http://localhost:11235", "Server address to connect to.")

// start as a zombie (defaults to human)
var zombie = flag.Bool("zombie", false, "Runs the example as a zombie.")

func main() {
	// eg. go run example.go -device COM2 -zombie
	flag.Parse()

	game.RegisterHuman(clock)
	// game.RegisterZombie(clock)

	err := game.Start("bob", *zombie, *device, *server)
	if err != nil {
		log.Fatal(err)
	}
}

// clock rotates 6 degrees every second
func clock(robot game.Robot) {
	tick := time.Tick(1 * time.Second)
	go func() {
		heading := 0 // 0-359 degrees
		speed := 1   // (stopped) 0-255 (fast)
		for {
			select {
			case <-tick: // Every tick
				robot.Walk(speed, heading)
				heading = (heading + 6) % 360
			case event, ok := <-robot.Events:
				if !ok {
					// Channel closed, I died.
					return
				}
				// Collision
				// Y Axis runs forwards/backwards (head on collisions)
				// positive values are the front (Y) & right (X)
				log.Printf("Event %+v\n.", event)

				// Turn around
				heading = (heading + 180) % 360
				robot.Walk(speed, heading)
			}
		}
	}()
}
