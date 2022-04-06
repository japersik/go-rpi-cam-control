package servoController

import (
	"github.com/stianeikeland/go-rpio/v4"
)

type sg90 struct {
	gpio    rpio.Pin
	inverse bool
	current int
}

const (
	leftLim  = 0
	rightLim = 180
	low      = 2000
	high     = 11000
)

func newServo(pin uint8, inverse bool) *sg90 {
	servo := &sg90{
		gpio:    rpio.Pin(pin),
		inverse: inverse,
		current: 90,
	}
	servo.gpio.Mode(rpio.Pwm)
	servo.gpio.Freq(4095000)
	return servo
}
func (s *sg90) Move(pos int) error {
	return s.MoveAbs(s.current + pos)
}

func (s *sg90) MoveAbs(pos int) error {
	s.current = pos
	if s.current > rightLim {
		s.current = rightLim
	} else if s.current < leftLim {
		s.current = leftLim
	}
	temp := s.current
	if s.inverse {
		temp = 180 - temp
	}
	s.gpio.DutyCycleWithPwmMode(uint32(low+(high-low)*(temp)/180), 81900, rpio.MarkSpace)
	return nil
}

func (s *sg90) CurrPosition() int {
	return s.current
}
