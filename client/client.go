// Package client handles communication with a go-zombie server gameroom.
package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/edmontongo/go-zombies/server/room"
	"github.com/edmontongo/gobot/platforms/sphero"
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

func (c *Client) Collide(data sphero.Collision) (room.Role, error) {
	json, err := wrap(data)
	if err != nil {
		return room.Unknown, err
	}
	request := fmt.Sprintf("%s/collision?id=%d&data=%s", c.roomUrl, c.id, json)
	var collision collisionResponse

	if err := getResponse(request, &collision); err != nil {
		return room.Unknown, err
	}
	log.Println("Hit a", collision.Hit)

	return room.ResolveRole(collision.Role), nil
}

func (c *Client) Close() error {
	request := fmt.Sprintf("%s/deregister?id=%d", c.roomUrl, c.id)
	var response deregisterResponse
	return getResponse(request, &response)
}

func wrap(i interface{}) (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
