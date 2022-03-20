package cameraController

type Camera interface {
	TakePhoto() (info PhotoInfo, err error)
	DelPhoto(id int) error
	GetPhoto(id int) (PhotoInfo, error)
	GetPhotoNums() int
	GetPhotoNames(count, pageNumber int) ([]PhotoInfo, error)
}
