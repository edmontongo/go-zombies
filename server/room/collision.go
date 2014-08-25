package room

import (
	"fmt"
	"time"

	"github.com/edmontongo/gobot/platforms/sphero"
)

type Collision struct {
	Id         Id
	ServerTime time.Time
	*player
	Collision sphero.Collision
	response  chan<- Role
}

func (c Collision) String() string {
	return fmt.Sprintf("%v: %s was %s to (??) from %+v", c.ServerTime, c.player.name, c.player.Role, c.Collision)
}

func (c Collision) Front() bool {
	return c.Collision.Y > 0 && ((c.Collision.Axis & 0x2) == 0x2)
}

func (c Collision) Strong() bool {
	return c.Collision.YMagnitude > 40
}

func (c Collision) Greater(other Collision) bool {
	// if c.Collision ==
	return false
}
