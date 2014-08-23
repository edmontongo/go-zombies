package game

import "log"

func defaultHuman(robot Robot) {
	log.Println("I'm a human")

	for {
		_, ok := <-robot.Events
		if !ok {
			break
		}
	}

	log.Println("Human death!")
}

func defaultZombie(robot Robot) {
	log.Println("I'm a zombie!")

	for {
		_, ok := <-robot.Events
		if !ok {
			break
		}
	}

	log.Println("Zombie death!")
}
