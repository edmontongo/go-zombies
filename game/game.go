package game

import (
	"fmt"
	"log"

	"github.com/edmontongo/go-zombies/client"
	"github.com/edmontongo/go-zombies/server/room"
	"github.com/edmontongo/gobot"
	"github.com/edmontongo/gobot/platforms/sphero"
)

// Robot handler
type Robot struct {
	adaptor *sphero.SpheroAdaptor
	driver  driver

	zombieFn robotFn
	humanFn  robotFn

	Events chan Event
	Died   chan int

	client *client.Client

	// Track our game state
	Role room.Role
}

type Event struct{}
type robotFn func(Robot)

var robot = Robot{
	Events: make(chan Event, 10),
	Died:   make(chan int, 1),

	zombieFn: defaultZombie,
	humanFn:  defaultHuman,
}

// RegisterHuman adds your brains to the game.
func RegisterHuman(fn robotFn) {
	robot.humanFn = fn
}

// RegisterZombie adds your zombie brains to the game.
func RegisterZombie(fn robotFn) {
	robot.zombieFn = fn
}

// Walk like a zombie.
func (r Robot) Walk(speed uint8, heading int) {
	r.driver.Roll(speed, uint16((heading+720)%360))
}

// Start the game
func Start(name string, zombie bool, port string) error {
	// TODO: pass in server IP
	c, err := client.New(name, "http://localhost:11235", zombie)
	if err != nil {
		return err
	}
	robot.client = c

	if zombie {
		robot.Role = room.Zombie
	} else {
		robot.Role = room.Human
	}

	if port == "" {
		fmt.Printf("Welcome %s.\n", name)
		// robot.driver = &fakeSpheroDriver{}
		fakeWork()
	} else {
		// TODO: err handling
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
	robot.driver.ConfigureCollisionDetectionRaw(0x10, 0x50, 0x10, 0x50, 0x60)

	// TODO: only if not a fakeSphero
	gobot.On(robot.driver.Event("collision"), func(data interface{}) {
		onCollission(data)
	})
	callUserCode()
}

func fakeWork() {
	// TODO: random events?
	callUserCode()
}

func onCollission(data interface{}) {
	fmt.Printf("Collision Detected! %+v\n", data)
	role, err := robot.client.Collide()
	if err != nil {
		log.Printf("Unexpected error during collision: %s", err)
		return
	}
	if role == room.Zombie {
		robot.driver.SetRGB(255, 0, 0)
	} else {
		robot.driver.SetRGB(0, 0, 255)
	}
	if robot.Role != role {
		robot.Died <- 1
		robot.Died = make(chan int, 1) // make a new channel for the new controlling function
		if role == room.Zombie {
			go robot.zombieFn(robot)
		} else {
			go robot.humanFn(robot)
		}
	}
	robot.Role = role
	robot.Events <- Event{}
}

func callUserCode() {
	// TODO: handle switching roles and all that
	if robot.Role == room.Zombie {
		robot.zombieFn(robot)
	} else {
		robot.humanFn(robot)
	}
}
