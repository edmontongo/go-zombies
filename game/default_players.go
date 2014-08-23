package game

import "log"

func defaultHuman(robot Robot) {
	log.Println("I'm a human")
	<-robot.Died
	log.Println("Human death!")
}

func defaultZombie(robot Robot) {
	log.Println("I'm a zombie!")
	<-robot.Died
	log.Println("Zombie death!")
}
