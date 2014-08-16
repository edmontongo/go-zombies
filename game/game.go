package game

import "fmt"

// import "github.com/edmontongo/go-zombies/client"

// Robot handler
type Robot struct {
	driver   driver
	zombieFn robotFn
	humanFn  robotFn
}

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
	if port == "" {
		fmt.Printf("Welcome %s.\n", name)
		robot.driver = &fakeSpheroDriver{}
	} else {
		// gobot := gobot.NewGobot()
		// adaptor := sphero.NewSpheroAdaptor(name, port)
		// robot.driver := sphero.NewSpheroDriver(adaptor, name)
	}

	// TODO: see which is registered, start only one
	robot.zombieFn(robot)

	return nil
}
