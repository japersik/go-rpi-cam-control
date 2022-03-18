package cameraController

import "io/fs"

type Camera interface {
	TakePhoto() (id int, err error)
	DelPhoto(id int) error
	GetPhoto(id int) (fs.FileInfo, error)
	GetPhotoNums() int
	GetPhotoNames(count, pageNumber int) ([]fs.FileInfo, error)
}
