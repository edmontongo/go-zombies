// Package room provides a Room in which multiple players enact a zombie simulation.
//
// The Room keeps track of the entire simulation state, requiring only collisions
// between players to be reported.
package room

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Room struct {
	players map[Id]*player
}

// Zombies returns a count of zombies in the room.
func (r *Room) Zombies() int {
	return r.count(Zombie)
}

// Humans returns a count of zombies in the room.
func (r *Room) Humans() int {
	return r.count(Human)
}

func (r *Room) count(role Role) int {
	count := 0
	for _, p := range r.players {
		if p.Role == role {
			count++
		}
	}
	return count
}

// Collide accepts two player ids and determines what there roles are after an interaction.
// If an error is returned than roles remain unchanged.
func (r *Room) Collide(id1, id2 Id) (r1, r2 Role, err error) {
	if id1 == id2 {
		err = fmt.Errorf("IDs %v and %v are equal", id1, id2)
		return
	}

	p1, err := r.player(id1)
	if err != nil {
		return
	}
	p2, err := r.player(id2)
	if err != nil {
		return
	}

	// Switch players for now, better math will be implemented later.
	p2.Role, p1.Role = p1.Role, p2.Role

	return p1.Role, p2.Role, nil
}

func (r *Room) player(id Id) (*player, error) {
	if p1, ok := r.players[id]; ok {
		return p1, nil
	}

	return nil, fmt.Errorf("player %v has not been registered", id)
}

// AddPlayer returns the unique player id after a player has been registered to the room.
func (r *Room) AddPlayer(name string, role Role, ip net.IP) Id {
	id := Id(0)
	for id == 0 {
		id = Id(rand.Int())
		if _, ok := r.players[id]; ok {
			id = 0
		}
	}

	r.players[id] = &player{name, id, role, ip, time.Now()}

	return id
}
