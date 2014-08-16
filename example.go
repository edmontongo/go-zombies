package main

import (
	"bufio"
	"log"
	"os"
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
				zombie.Walk(1, uint16(a%360))
				a += 15
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
	err := game.Start("bob", "")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
