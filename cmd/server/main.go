package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController/rpiCamera"
	"github.com/japersik/go-rpi-cam-control/internal/app/moveController/ServoController"
	"github.com/japersik/go-rpi-cam-control/internal/app/server"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configs-path", "configs/serverConf.toml", "path to configs file")

}
func main() {
	flag.Parse()

	config := server.NewConfig()
	config2 := ServoController.NewConfig()
	config3 := rpiCamera.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	_, err = toml.DecodeFile(configPath, config2)
	_, err = toml.DecodeFile(configPath, config3)
	if err != nil {
		log.Fatal(err)
	}
	servo, _ := ServoController.NewServoController(config2)
	camera := rpiCamera.NewRpiCam(config3)
	fmt.Println(config2)
	if err := server.Start(config, servo, camera); err != nil {
		log.Fatal(err)
	}
}
