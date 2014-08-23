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

	client *client.Client

	// Track our game state
	Role room.Role
}

// Event is a game event, like a collision.
type Event struct{}
type robotFn func(Robot)

var robot = Robot{
	Events: make(chan Event, 10),

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
	// threshold at speed of 0, threshold at maximum speed
	robot.driver.ConfigureCollisionDetectionRaw(0x40, 0x40, 0x50, 0x50, 0x60)
	robot.driver.SetBackLED(0xff)

	// TODO: only if not a fakeSphero
	gobot.On(robot.driver.Event("collision"), func(data interface{}) {
		collision, _ := data.(sphero.Collision)
		onCollission(collision)
	})
	robot.setColor(robot.Role)
	lauchPlayerCode()
}

func fakeWork() {
	// TODO: random events?
	lauchPlayerCode()
}

func onCollission(collision sphero.Collision) {
	fmt.Printf("Collision Detected! %+v\n", collision)
	// Y Axis runs forwards/backwards (head on collisions)
	// positive values are the front (Y) & right (X)
	role, err := robot.client.Collide()
	if err != nil {
		log.Printf("Unexpected error during collision: %s", err)
		return
	}

	robot.setColor(role)
	robot.Events <- Event{}

	if robot.Role != role {
		// restart the event loop
		close(robot.Events)
		robot.Events = make(chan Event, 10)
		robot.Role = role
		lauchPlayerCode()
	}
}

func lauchPlayerCode() {
	// make new channels for the new controlling function

	if robot.Role == room.Zombie {
		go robot.zombieFn(robot)
	} else {
		go robot.humanFn(robot)
	}
}

func (r *Robot) setColor(role room.Role) {
	if role == room.Zombie {
		robot.driver.SetRGB(255, 0, 0)
	} else {
		robot.driver.SetRGB(0, 0, 255)
	}
}
