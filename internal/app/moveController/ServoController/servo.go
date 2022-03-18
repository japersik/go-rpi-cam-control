package ServoController

import "github.com/stianeikeland/go-rpio/v4"

type sg90 struct {
	gpio    rpio.Pin
	inverse bool
	current int
}

const (
	leftLim  = -10
	rightLim = 190
)

func newServo(pin uint8, inverse bool) *sg90 {
	servo := &sg90{
		gpio:    rpio.Pin(pin),
		inverse: inverse,
		current: 90,
	}
	servo.gpio.Output()
	return servo
}
func (s *sg90) Move(pos int) error {
	return s.MoveAbs(s.current + pos)
}

func (s *sg90) MoveAbs(pos int) error {
	if pos > rightLim {
		pos = rightLim
	} else if pos < leftLim {
		pos = leftLim
	}
	s.current = pos
	s.gpio.DutyCycleWithPwmMode(0, 200, rpio.MarkSpace)
	return nil
}
func (s *sg90) CurrPosition() int {
	return s.current
}
