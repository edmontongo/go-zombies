// Package client handles communication with a go-zombie server gameroom.
package client

import (
	"fmt"
	"log"

	"github.com/edmontongo/go-zombies/server/room"
)

type Client struct {
	roomUrl string
	id      room.Id
}

func New(name, url string, zombie bool) (*Client, error) {
	request := fmt.Sprintf("%s/register?name=%s", url, name)
	if zombie {
		request += "&role=zombie"
	}

	var register registerResponse
	if err := getResponse(request, &register); err != nil {
		return nil, err
	}

	c := Client{
		id:      room.Id(register.PlayerId),
		roomUrl: url,
	}

	return &c, nil
}

func (c *Client) Collide() (room.Role, error) {
	request := fmt.Sprintf("%s/collision?id=%d", c.roomUrl, c.id)
	var collision collisionResponse

	if err := getResponse(request, &collision); err != nil {
		return room.Unknown, err
	}

	log.Println("Hit a", collision.Hit)

	return room.ResolveRole(collision.Role), nil
}
