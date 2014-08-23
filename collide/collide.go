package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edmontongo/go-zombies/game"
)

var port = flag.String("port", "/dev/tty.Sphero-OBY-RN-SPP", "Port to the Sphero")

// true to start as a zombie
var zombie = flag.Bool("zombie", false, "Runs the example as a zombie")

func humanRoller(human game.Robot) {
	go func() {
		b := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf(":> ")
			line, err := b.ReadString('\n')
			if err != nil {
				log.Fatal(err.Error())
			}

			switch line[0] {
			case 's':
				clear(human.Events)
				var speed uint8
				var heading int
				fmt.Sscanf(line[2:], "%d %d", &speed, &heading)
				human.Walk(speed, heading)
				c := time.After(5 * time.Second)
				select {
				case <-c:
					human.Walk(0, heading)
					log.Println("Timedout")
				case event, ok := <-human.Events:
					if !ok {
						return
					}
					log.Printf("Event %v\n.", event)
					human.Walk(0, heading)
				}

				clear(human.Events)

			case 'c':
				var xThreshold, xSpeed, yThreshold, ySpeed uint8

				fmt.Sscanf(line[2:], "%d %d %d %d", &xThreshold, &xSpeed, &yThreshold, &ySpeed)
				human.Driver.ConfigureCollisionDetectionRaw(xThreshold, xSpeed, yThreshold, ySpeed, 25)
				log.Println("Set collision detection", xThreshold, xSpeed, yThreshold, ySpeed)
			}

		}
	}()
}

func clear(c chan game.Event) {
	for {
		select {
		case event, ok := <-c:
			if !ok {
				return
			}
			log.Printf("Event %v\n.", event)
		default:
			return
		}
	}
}

func main() {
	flag.Parse()

	game.RegisterHuman(humanRoller)
	game.RegisterZombie(humanRoller)

	err := game.Start("bob", *zombie, *port)
	if err != nil {
		log.Fatal(err)
	}
}
