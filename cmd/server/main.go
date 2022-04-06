package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController/rpiCamera"
	"github.com/japersik/go-rpi-cam-control/internal/app/moveController/servoController"
	"github.com/japersik/go-rpi-cam-control/internal/app/server"
	"github.com/japersik/go-rpi-cam-control/internal/app/service"
	"github.com/japersik/go-rpi-cam-control/internal/app/userStore/jsonStore"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configs-path", "configs/serverConf.toml", "path to configs file")

}

type Configs struct {
	ServerConfig    server.Config          `toml:"server_config"`
	CameraConfig    rpiCamera.Config       `toml:"camera_config"`
	MoveConfig      servoController.Config `toml:"servo_config"`
	UserStoreConfig jsonStore.Config       `toml:"user_store"`
}

func main() {
	flag.Parse()
	config := new(Configs)
	_, err := toml.DecodeFile(configPath, config)
	tomlWriter := toml.NewEncoder(os.Stdout)
	tomlWriter.Encode(config)

	if err != nil {
		log.Fatal(err)
	}
	servo, _ := servoController.NewServoController(&config.MoveConfig)
	camera := rpiCamera.NewRpiCam(&config.CameraConfig)
	userData := jsonStore.NewUserStore(&config.UserStoreConfig)
	rpiService := service.Service{userData, servo, camera}
	if err := server.Start(&config.ServerConfig, rpiService); err != nil {
		log.Fatal(err)
	}
}
