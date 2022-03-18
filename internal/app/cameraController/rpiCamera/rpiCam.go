package rpiCamera

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strconv"
	"strings"
)

type rpiCam struct {
	photoDir string
	idToPath map[int]fs.FileInfo
}

func NewRpiCam(conf *Config) *rpiCam {
	path := "private/static/img"
	var mapa map[int]fs.FileInfo
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	if files != nil {
		mapa = make(map[int]fs.FileInfo, len(files))
		for i := 0; i < len(files); i++ {
			if strings.HasSuffix(files[i].Name(), ".png") {
				mapa[len(mapa)+1] = files[i]
				fmt.Println(files[i].Name())
			}
		}
	} else {
		mapa = make(map[int]fs.FileInfo, 20)
	}
	cam := &rpiCam{
		photoDir: "path",
		idToPath: mapa,
	}
	fmt.Println(cam.idToPath)
	return cam
}

func (r *rpiCam) TakePhoto() (id int, err error) {
	fmt.Println("newPhoto")
	return 0, nil
}
func (r *rpiCam) DelPhoto(id int) (err error) {
	fmt.Println("delPhoto", id)
	return nil
}
func (r *rpiCam) GetPhoto(id int) (fs.FileInfo, error) {
	if name, ok := r.idToPath[id]; ok {
		return name, nil
	}
	return nil, errors.New("id not found")
}
func (r *rpiCam) GetPhotoNums() int {
	return len(r.idToPath)
}
func (r *rpiCam) GetPhotoNames(count, pageNumber int) ([]fs.FileInfo, error) {
	if count*(pageNumber-1) > r.GetPhotoNums() {
		return nil, errors.New("Maximum page: " + strconv.Itoa(r.GetPhotoNums()))
	}
	names := make([]fs.FileInfo, 0, count)
	counter := 0
	fmt.Println(r.idToPath)
	for _, s := range r.idToPath {
		counter++
		if counter < (count * (pageNumber - 1)) {
			continue
		}
		if counter > (count * (pageNumber + 1)) {
			break
		}
		names = append(names, s)
	}
	return names, nil
}
