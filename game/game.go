package game

import (
	"fmt"
	"log"

	"github.com/edmontongo/go-zombies/client"
	"github.com/edmontongo/go-zombies/server/room"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

// Robot handler
type Robot struct {
	adaptor *sphero.SpheroAdaptor
	driver  driver

	zombieFn robotFn
	humanFn  robotFn

	Events chan Event

	client *client.Client

	// Track our game state
	Role room.Role
}

type Event struct{}
type robotFn func(Robot)

var robot Robot

// RegisterZombie does stuff
func RegisterZombie(fn robotFn) {
	robot.zombieFn = fn
}

// RegisterHuman does stuff
func RegisterHuman(fn robotFn) {
	robot.humanFn = fn
}

// Walk like a zombie.
func (r Robot) Walk(speed uint8, heading uint16) {
	r.driver.Roll(speed, heading)
}

// Start the game
func Start(name string, zombie bool, port string) error {
	c, err := client.New(name, "http://localhost:11235", false)
	if err != nil {
		return err
	}
	robot.client = c
	robot.Events = make(chan Event, 10)
	if zombie {
		robot.Role = room.Zombie
	} else {
		robot.Role = room.Human
	}

	if port == "" {
		fmt.Printf("Welcome %s.\n", name)
		// robot.driver = &fakeSpheroDriver{}
		work()
	} else {
		bot := gobot.NewGobot()
		robot.adaptor = sphero.NewSpheroAdaptor(name, port)
		robot.driver = sphero.NewSpheroDriver(robot.adaptor, name)

		sphero := gobot.NewRobot(name,
			[]gobot.Connection{robot.adaptor},
			[]gobot.Device{robot.driver},
			work,
		)
		bot.AddRobot(sphero)
		bot.Start()
	}

	return nil
}

func work() {
	// TODO: only if not a fakeSphero
	gobot.On(robot.driver.Event("collision"), func(data interface{}) {
		fmt.Printf("Collision Detected! %+v\n", data)
		role, err := robot.client.Collide()
		if err != nil {
			log.Printf("Unexpected error during collision: %s", err.Error())
			return
		}
		if role == room.Zombie {
			robot.driver.SetRGB(255, 0, 0)
		} else {
			robot.driver.SetRGB(0, 0, 255)
		}
		robot.Role = role
		robot.Events <- Event{}
	})

	// TODO: see which is registered, start only one
	robot.zombieFn(robot)
}
