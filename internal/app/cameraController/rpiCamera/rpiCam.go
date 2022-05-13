package rpiCamera

import (
	"errors"
	"fmt"
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type rpiCam struct {
	photoDir string
	photos   []cameraController.PhotoInfo
}

func NewRpiCam(conf *Config) *rpiCam {
	path := "private/static/img/"
	var phArr []cameraController.PhotoInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	if files != nil {
		phArr = make([]cameraController.PhotoInfo, 0, len(files))
		for i := 0; i < len(files); i++ {
			if strings.HasSuffix(files[i].Name(), ".jpeg") || strings.HasSuffix(files[i].Name(), ".png") {
				phArr = append(phArr, cameraController.PhotoInfo{Name: files[i].Name(), Size: files[i].Size(), CreationTime: files[i].ModTime()})
			}
		}
	} else {
		phArr = make([]cameraController.PhotoInfo, 0, 40)
	}
	cam := &rpiCam{
		photoDir: path,
		photos:   phArr,
	}
	return cam
}

func (r *rpiCam) TakePhoto() (cameraController.PhotoInfo, error) {
	t := time.Now()
	name := r.photoDir + t.Format("2006-01-02-150405") + strconv.Itoa(len(r.photos)+1) + ".jpeg"
	_, err := cmdTakePhoto(name)
	if err != nil {
		return cameraController.PhotoInfo{}, err
	}
	fs, err := os.Stat(name)
	if err != nil {
		return cameraController.PhotoInfo{}, err
	}
	r.photos = append(r.photos, cameraController.PhotoInfo{Name: fs.Name(), Size: fs.Size(), CreationTime: fs.ModTime()})

	return r.photos[len(r.photos)-1], nil
}
func (r *rpiCam) DelPhoto(id int) (err error) {
	for i := 0; i < len(r.photos); i++ {
		fmt.Printf("%d-", i)
		if r.photos[i].Id == id {

			fmt.Println(r.photoDir + r.photos[i].Name)
			err := os.Remove(r.photoDir + r.photos[i].Name)
			if err != nil {
				return err
			}
			r.photos = append(r.photos[0:i], r.photos[i+1:len(r.photos)]...)
			return nil
		}
	}
	return nil
}
func (r *rpiCam) GetPhoto(id int) (cameraController.PhotoInfo, error) {
	if len(r.photos) > id {
		return r.photos[id], nil
	}
	return cameraController.PhotoInfo{}, errors.New("id not found")
}

func cmdTakePhoto(name string) ([]byte, error) {
	app := "raspistill"

	arg0 := "-o"
	arg1 := name
	arg2 := "-q"
	arg3 := "20"
	arg4 := "-t"
	arg5 := "150"
	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5)
	return cmd.Output()
}
func (r *rpiCam) GetPhotoNums() int {
	return len(r.photos) - 1
}

func (r *rpiCam) GetPhotoNames(count, pageNumber int) ([]cameraController.PhotoInfo, error) {
	if count*(pageNumber-1) > r.GetPhotoNums() {
		return nil, errors.New("Maximum page: " + strconv.Itoa(r.GetPhotoNums()))
	}
	names := make([]cameraController.PhotoInfo, 0, count)
	left := len(r.photos) - (pageNumber-1)*count - 1
	right := left - count
	if right < 0 {
		right = 0
	}
	fmt.Println(right, left)
	for ; left > right; left-- {
		r.photos[left].Id = left
		names = append(names, r.photos[left])
	}
	return names, nil
}
