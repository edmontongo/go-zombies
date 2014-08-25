package room

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Role uint

const (
	Unknown = Role(iota)
	Zombie
	Human
	Wall
	lastRole
)

var roleStrings = []string{
	"Unknown",
	"Zombie",
	"Human",
	"Wall",
}

func (r Role) String() string {
	if r > lastRole {
		r = 0
	}
	return roleStrings[r]
}

func ResolveRole(s string) Role {
	for i, roleString := range roleStrings {
		if s == roleString {
			return Role(i)
		}
	}
	return Unknown
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

func (p player) Description() string {
	return fmt.Sprintf("%s: %s", p.name, p.Role)
}
