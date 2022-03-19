package ServoController

import (
	"github.com/stianeikeland/go-rpio/v4"
)

//func main() {
//
//
//
//	pin := rpio.Pin(13)
//	pin.Mode(rpio.Pwm)
//	pin.Freq(5000)
//	pin.DutyCycleWithPwmMode(0, 200, rpio.MarkSpace)
//	// the LED will be blinking at 2000Hz
//	// (source frequency divided by cycle length => 64000/32 = 2000)
//
//	// five times smoothly fade in and out
//	for i := 5; i < 100; i+=5 {
//		fmt.Println(i)
//		pin.DutyCycleWithPwmMode(uint32(i),200, rpio.MarkSpace)
//		time.Sleep(time.Second)
//	}
//
//	pin.DutyCycleWithPwmMode(0, 200, rpio.Balanced)
//}

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
