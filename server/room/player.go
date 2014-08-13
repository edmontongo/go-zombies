package room

import (
	"net"
	"strconv"
	"time"
)

type Role uint

const (
	Unknown = Role(iota)
	Zombie
	Human
	lastRole
)

var roleStrings = []string{
	"Unknown",
	"Zombie",
	"Human",
}

func (r Role) String() string {
	if r > lastRole {
		r = 0
	}
	return roleStrings[r]
}

// Id is a unique (per room) identifier for communicating player ids with the Room.
type Id int

func IdFromString(s string) (Id, error) {
	id, err := strconv.Atoi(s)
	return Id(id), err
}

type player struct {
	name string
	Id
	Role
	ip     net.IP
	joined time.Time
}
