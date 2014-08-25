// Package room provides a Room in which multiple players enact a zombie simulation.
//
// The Room keeps track of the entire simulation state, requiring only collisions
// between players to be reported.
package room

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"
)

type Room struct {
	players          map[Id]*player
	collisionQueue   chan<- *Collision
	recentCollisions []*Collision
}

func New() Room {
	c := make(chan *Collision)
	r := Room{map[Id]*player{}, c, make([]*Collision, 10)}
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
func (r *Room) collide(c1, c2 *Collision) (r1, r2 Role) {
	if c2.Role == c1.Role {
		return c1.Role, c2.Role
	}

	// Switch players for now, better math will be implemented later.
	if c1.Collision.Speed > c2.Collision.Speed {
		c2.Role = c1.Role
	} else {
		c1.Role = c2.Role
	}

	return c1.Role, c2.Role
}

// Collision checks if the given id was involved in a collision with anyone else. An error is returned if the player wasn't registered to the room.
func (r *Room) Collision(c Collision) (newRole, hit Role, err error) {
	c.player, err = r.player(c.Id)
	if err != nil {
		return Invalid, Invalid, err
	}

	if !c.Valid() {
		return c.Role, Unknown, nil
	}

	result := make(chan Role)
	c.response = result
	r.collisionQueue <- &c

	return c.player.Role, <-result, nil
}

func (r *Room) collisionManager(c <-chan *Collision) {
	for c1 := range c {
	top:
		t := time.After(400 * time.Millisecond)
		select {
		case c2 := <-c:
			if c1.Attack() == c2.Attack() {
				c1.response <- Unknown
				c1 = c2
				goto top
			}

			r.recentCollisions = append(r.recentCollisions[2:], c1, c2)
			log.Println("What??!?!")
			oldp1, oldp2 := c1.Role, c2.Role
			r.collide(c1, c2)
			c2.response <- oldp1
			c1.response <- oldp2
		case _ = <-t:
			c1.response <- Unknown
		}
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
	rand.Seed(time.Now().UnixNano())
	for id == 0 {
		id = Id(rand.Int())
		if _, ok := r.players[id]; ok {
			id = 0
		}
	}

	r.players[id] = &player{name, id, role, ip, time.Now()}
	log.Printf("AddPlayer %v\n", id)
	return id
}

// RemovePlayer removes a player from the room.
func (r *Room) RemovePlayer(id Id) {
	log.Printf("RemovePlayer %v\n", id)
	delete(r.players, id)
}

func (r *Room) Recent() []*Collision {
	ret := make([]*Collision, len(r.recentCollisions))
	copy(ret, r.recentCollisions)
	return ret
}

type playerList []*player

func (p playerList) Len() int           { return len(p) }
func (p playerList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p playerList) Less(i, j int) bool { return p[i].name < p[j].name }

// Players returns a sorted list of players
func (r *Room) Players() []*player {
	list := []*player{}
	for _, p := range r.players {
		list = append(list, p)
	}
	sort.Sort(playerList(list))
	return list
}
