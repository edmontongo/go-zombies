package game

import (
	"fmt"

	"github.com/hybridgroup/gobot"
)

type driver interface {
	gobot.DriverInterface
	SetHeading(heading uint16)
	Roll(speed uint8, heading uint16)
	// Halt() bool
	SetBackLED(level uint8)
	SetRGB(r uint8, g uint8, b uint8)
}

type fakeSpheroDriver struct{}

func (s *fakeSpheroDriver) SetHeading(heading uint16) {
	fmt.Printf("Heading set to %d.\n", heading)
}

func (s *fakeSpheroDriver) Roll(speed uint8, heading uint16) {
	fmt.Printf("Rolling, rolling rolling... speed %d, heading %d.\n", speed, heading)
}

func (s *fakeSpheroDriver) Halt() bool {
	fmt.Println("Halt!")
	return true
}

func (s *fakeSpheroDriver) SetBackLED(level uint8) {
	fmt.Println("Light the way: %d.", level)
}

func (s *fakeSpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	fmt.Println("Rainbows: #%2x%2x%2x", r, g, b)
}
