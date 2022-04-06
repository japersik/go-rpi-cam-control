package servoController

import (
	"github.com/stianeikeland/go-rpio/v4"
)

type GPIOServoController struct {
	servoX *sg90
	servoY *sg90
	opened bool
}

// Config ...
type Config struct {
	GpioX    int  `toml:"gpio_x"`
	GpioY    int  `toml:"gpio_y"`
	InverseX bool `toml:"inverse_x"`
	InverseY bool `toml:"inverse_y"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		GpioX:    19,
		GpioY:    12,
		InverseX: false,
		InverseY: false,
	}
}
func NewServoController(config *Config) (*GPIOServoController, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	controller := &GPIOServoController{
		servoX: newServo(uint8(config.GpioX), config.InverseX),
		servoY: newServo(uint8(config.GpioY), config.InverseX),
	}
	return controller, nil
}

func (s *GPIOServoController) Close() {
	rpio.Close()
}

func (s *GPIOServoController) MoveX(pos int) error {
	return s.servoX.Move(pos)
}

func (s *GPIOServoController) MoveY(pos int) error {
	return s.servoY.Move(pos)
}
func (s *GPIOServoController) MoveXAbs(pos int) error {
	return s.servoX.MoveAbs(pos)
}
func (s *GPIOServoController) MoveYAbs(pos int) error {
	return s.servoY.MoveAbs(pos)
}
func (s *GPIOServoController) CurrPosition() (x, y int) {
	return s.servoX.CurrPosition(), s.servoY.CurrPosition()
}
