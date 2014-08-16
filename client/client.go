// Package client handles communication with a go-zombie server gameroom.
package client

import (
	"fmt"

	"github.com/edmontongo/go-zombies/server/room"
)

type Client struct {
	roomUrl string
	id      room.Id
}

func New(name, url string) (*Client, error) {
	var register registerResponse
	if err := getResponse(fmt.Sprintf("%s/register?name=%s", url, name), &register); err != nil {
		return nil, err
	}

	c := Client{
		id:      room.Id(register.PlayerId),
		roomUrl: url,
	}

	return &c, nil
}
