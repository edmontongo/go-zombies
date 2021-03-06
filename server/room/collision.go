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

func (c Collision) Valid() bool {
	return c.Vulnerable() || c.Attack()
}

func (c Collision) Attack() bool {
	return c.Front() && c.Strong() && c.Collision.Speed > 60
}

func (c Collision) Vulnerable() bool {
	return c.Collision.Speed < 20 && (c.Collision.XMagnitude+c.Collision.YMagnitude) > 80
}

func (c Collision) Strong() bool {
	return c.Collision.YMagnitude > 60
}

func (c Collision) Greater(other Collision) bool {
	// if c.Collision ==
	return false
}
