package game

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

// import "github.com/edmontongo/go-zombies/client"

// Robot handler
type Robot struct {
	adaptor *sphero.SpheroAdaptor
	driver  driver

	zombieFn robotFn
	humanFn  robotFn

	Events chan Event
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
func Start(name string, port string) error {
	robot.Events = make(chan Event, 10)

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
		robot.Events <- Event{}
	})

	// TODO: see which is registered, start only one
	robot.zombieFn(robot)
}
