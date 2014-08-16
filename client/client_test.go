package client

// warning: Tests expect a server running at TestURL

import (
	"testing"

	"github.com/edmontongo/go-zombies/server/room"
	"github.com/stretchr/testify/assert"
)

const testURL = "http://localhost:11235"

func TestRegistration(t *testing.T) {
	c, err := New("Bob", testURL, false)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, c.id)
}

func TestZombieRegistration(t *testing.T) {
	c, err := New("Bob", testURL, true)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, c.id)
}

func TestRegsistrationNoName(t *testing.T) {
	c, err := New("", testURL, false)
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestCollision(t *testing.T) {
	bob, err := New("Bob", testURL, false)
	assert.NoError(t, err)
	alice, err := New("Alice", testURL, true)
	assert.NoError(t, err)

	wait := make(chan int, 1)
	var bobRole room.Role
	go func() {
		var err error
		bobRole, err = bob.Collide()
		assert.NoError(t, err, "Unexpected Bob error")
		wait <- 1
	}()

	aliceRole, err := alice.Collide()
	assert.NoError(t, err, "Unexpected Alice error")
	<-wait
	assert.Equal(t, aliceRole, bobRole)
}
