package game

import (
	"log"

	"github.com/edmontongo/go-zombies/client"
	"github.com/edmontongo/go-zombies/server/room"
	"github.com/edmontongo/gobot"
	"github.com/edmontongo/gobot/platforms/sphero"
)

// Robot that can walk around.
type Robot struct {
	adaptor *sphero.SpheroAdaptor
	driver  *sphero.SpheroDriver

	zombieFn robotFn
	humanFn  robotFn

	Events chan Event

	client *client.Client

	// Track our game state
	Role room.Role

	humanColor *Color
	// mHumanColor sync.RWMutex
}

// Event is a game event, like a collision.
type Event struct {
	Collision             sphero.Collision
	WasRole, NewRole, Hit room.Role
}

type robotFn func(Robot)

// Color (red, green, blue)
type Color struct {
	Red, Green, Blue uint8
}

// one global robot
var robot = Robot{
	Events: make(chan Event, 10),

	zombieFn:   defaultZombie,
	humanFn:    defaultHuman,
	humanColor: &Color{0, 0, 255},
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
func (r *Robot) Walk(speed, heading int) {
	heading = (heading + 720) % 360
	r.driver.Roll(uint8(speed), uint16(heading))
}

// SetReferenceHeading calibrates a heading (0-359).
func (r *Robot) SetReferenceHeading(heading int) {
	heading = (heading + 720) % 360
	r.driver.SetHeading(uint16(heading))
}

// SetHumanColor sets the red, green, blue color value.
func (r *Robot) SetHumanColor(red, green, blue uint8) {
	// r.mHumanColor.Lock()
	r.humanColor.Red = red
	r.humanColor.Green = green
	r.humanColor.Blue = blue
	// r.mHumanColor.Unlock()
	r.setColor(r.Role)
	log.Printf("set: %v\n", r.humanColor)
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
	// threshold at speed of 0, add value at maximum speed
	robot.driver.ConfigureCollisionDetectionRaw(20, 0, 20, 0, 50)
	robot.driver.SetBackLED(0xff)

	gobot.On(robot.driver.Event("collision"), func(data interface{}) {
		collision, _ := data.(sphero.Collision)
		onCollission(collision)
	})
	robot.setColor(robot.Role)
	lauchPlayerCode()
}

func onCollission(collision sphero.Collision) {
	// Y Axis runs forwards/backwards (head on collisions)
	// positive values are the front (Y) & right (X)
	role, hitRole, err := robot.client.Collide(collision)
	if err != nil {
		log.Printf("Unexpected error sending collision to server: %s", err)
		return
	}

	robot.setColor(role)
	event := Event{
		Collision: collision,
		WasRole:   robot.Role,
		NewRole:   role,
		Hit:       hitRole,
	}
	robot.Events <- event

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
		// r.mHumanColor.RLock()
		// log.Println(r.humanColor)
		robot.driver.SetRGB(r.humanColor.Red, r.humanColor.Green, r.humanColor.Blue)
		// r.mHumanColor.RUnlock()
	}
}
