package cameraController

import "time"

type PhotoInfo struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	CreationTime time.Time `json:"creationTime"`
}
