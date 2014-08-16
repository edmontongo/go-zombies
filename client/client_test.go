package client

// warning: Tests expect a server running at TestURL

import (
	"testing"

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
