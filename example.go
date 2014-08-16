package main

import (
	"log"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

/*
	Mac:
	Pair with a Bluetooth device, then specify a port in the form:
 		/dev/tty.Sphero-???-AMP-SPP
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
var port = "/dev/tty.Sphero-WOO-AMP-SPP"

func zombieTicker(zombie game.Robot) {
	c := time.Tick(1 * time.Second)
	go func() {
		a := 0
		for {
			select {
			case <-c:
				// do stuff
				zombie.Walk(0, uint16(a%360))
				a += 6
			case event, ok := <-zombie.Events:
				if !ok {
					return
				}
				log.Printf("Event %v\n.", event)
			}
		}
	}()
}

func main() {
	game.RegisterZombie(zombieTicker)
	// game.RegisterHuman(me)

	err := game.Start("bob", port)
	if err != nil {
		log.Fatal(err)
	}
}
