package service

import (
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController"
	"github.com/japersik/go-rpi-cam-control/internal/app/moveController"
	"github.com/japersik/go-rpi-cam-control/internal/app/userStore"
)

type Service struct {
	userStore.Authorization
	moveController.Mover
	cameraController.Camera
}
