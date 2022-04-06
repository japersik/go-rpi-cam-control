package server

import (
	"github.com/japersik/go-rpi-cam-control/internal/app/service"
	"net/http"
)

func Start(config *Config, service service.Service) error {

	s := newServer(config, service)
	return http.ListenAndServe(config.BindAddr, s)
}
