package server

import (
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController"
	"github.com/japersik/go-rpi-cam-control/internal/app/moveController"
	"net/http"
)

func Start(config *Config, mover moveController.Mover, camera cameraController.Camera) error {

	s := newServer(config, mover, camera)
	s.configureLogger(config.LogLevel)
	return http.ListenAndServe(config.BindAddr, s)
}
