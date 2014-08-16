package main

import (
	"log"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

func sleepyZombie(zombie game.Robot) {
	c := time.Tick(1000 * time.Millisecond)
	go func() {
		a := 0
		for {
			select {
			case <-c:
				// do stuff
				zombie.Walk(0, uint16(a%360))
				a += 6
				// case event, ok = <-zombie.Event:
				// 	if !ok {
				// 		return
				// 	}
				// collision
			}
		}
	}()
}

func main() {
	game.RegisterZombie(sleepyZombie)
	// game.RegisterHuman(me)

	// Mac:
	// /dev/tty.Sphero-WOO-AMP-SPP
	// Windows:
	// Linux:
	err := game.Start("bob", "/dev/tty.Sphero-WOO-AMP-SPP")
	if err != nil {
		log.Fatal(err)
	}

	// reader := bufio.NewReader(os.Stdin)
	// reader.ReadString('\n')
}
