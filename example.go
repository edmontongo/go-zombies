package main

import (
	"fmt"
	"log"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

func zombieClock(zombie game.Robot) {
	c := time.Tick(1000 * time.Millisecond)
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
				fmt.Printf("Event %v\n.", event)
			}
		}
	}()
}

func main() {
	game.RegisterZombie(zombieClock)
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
