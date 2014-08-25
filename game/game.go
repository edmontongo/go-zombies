package game

import (
	"log"

	"github.com/edmontongo/go-zombies/client"
	"github.com/edmontongo/go-zombies/server/room"
	"github.com/edmontongo/gobot"
	"github.com/edmontongo/gobot/platforms/sphero"
)

// Robot handler
type Robot struct {
	adaptor *sphero.SpheroAdaptor
	driver  *sphero.SpheroDriver

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
func (r Robot) Walk(speed, heading int) {
	r.driver.Roll(uint8(speed), uint16((heading+720)%360))
}

// SetReferenceHeading calibrates a heading (0-359).
func (r Robot) SetReferenceHeading(heading int) {
	r.driver.SetHeading(uint16((heading + 720) % 360))
}

// Start the game
func Start(name string, zombie bool, device string, server string) error {
	moreWork := func() {
		c, err := client.New(name, server, zombie)
		if err != nil {
			log.Fatal(err)
		}
		robot.client = c

		if zombie {
			robot.Role = room.Zombie
		} else {
			robot.Role = room.Human
		}

		work()
	}

	// TODO: err handling
	bot := gobot.NewGobot()
	robot.adaptor = sphero.NewSpheroAdaptor(name, device)
	robot.driver = sphero.NewSpheroDriver(robot.adaptor, name)

	sphero := gobot.NewRobot(name,
		[]gobot.Connection{robot.adaptor},
		[]gobot.Device{robot.driver},
		moreWork,
	)
	bot.AddRobot(sphero)

	bot.Start()
	robot.client.Close()

	return nil
}

func work() {
	// threshold at speed of 0, threshold at maximum speed
	robot.driver.ConfigureCollisionDetectionRaw(20, 0, 20, 0, 50)
	robot.driver.SetBackLED(0xff)

	// TODO: only if not a fakeSphero
	gobot.On(robot.driver.Event("collision"), func(data interface{}) {
		collision, _ := data.(sphero.Collision)
		onCollission(collision)
	})
	robot.setColor(robot.Role)
	lauchPlayerCode()
}

func onCollission(collision sphero.Collision) {
	log.Printf("Collision Detected! %+v\n", collision)
	// Y Axis runs forwards/backwards (head on collisions)
	// positive values are the front (Y) & right (X)
	role, err := robot.client.Collide(collision)
	if err != nil {
		log.Printf("Unexpected error sending collision to server: %s", err)
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
