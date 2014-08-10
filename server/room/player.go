package room

import (
	"net"
	"time"
)

type Role int

const (
	Unknown = Role(iota)
	Zombie
	Human
)

// Id is a unique (per room) identifier for communicating player ids with the Room.
type Id int

type player struct {
	name string
	Id
	Role
	net.IP
	joined time.Time
}
