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

type queuedCollision struct {
	*player
	response chan<- Role
}

type Room struct {
	players        map[Id]*player
	collisionQueue chan<- queuedCollision
}

func New() Room {
	c := make(chan queuedCollision)
	r := Room{map[Id]*player{}, c}
	go r.collisionManager(c)
	return r
}

func (r *Room) Close() {
	close(r.collisionQueue)
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

// Collide accepts two players and determines what there roles are after an interaction.
func (r *Room) collide(p1, p2 *player) (r1, r2 Role) {
	if p2.Role == p1.Role {
		return p1.Role, p2.Role
	}

	winner := rand.Float32() > 0.70
	// Switch players for now, better math will be implemented later.
	if p2.Role == Zombie {
		if winner {
			p1.Role = Zombie
		} else {
			p2.Role = Human
		}
	} else {
		if winner {
			p2.Role = Zombie
		} else {
			p1.Role = Human
		}
	}

	return p1.Role, p2.Role
}

/// Collision checks if the given id was involved in a collision with anyone else. An error is returned if the player wasn't registered to the room.
func (r *Room) Collision(id Id) (Role, error) {
	p, err := r.player(id)
	if err != nil {
		return p.Role, err
	}

	c := make(chan Role)
	r.collisionQueue <- queuedCollision{p, c}
	<-c

	return p.Role, nil
}

func (r *Room) collisionManager(c <-chan queuedCollision) {
	for p1 := range c {
		t := time.After(time.Second)
		select {
		case p2 := <-c:
			r.collide(p1.player, p2.player)
			p2.response <- p2.Role
		case _ = <-t:
		}
		p1.response <- p1.Role
	}
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
