package rpiCamera

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type rpiCam struct {
	photoDir  string
	idToPath  map[int]fs.FileInfo
	idCounter int
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
			if strings.HasSuffix(files[i].Name(), ".jpeg") {
				mapa[len(mapa)+1] = files[i]
				fmt.Println(files[i].Name())
			}
		}
	} else {
		mapa = make(map[int]fs.FileInfo, 20)
	}
	cam := &rpiCam{
		photoDir: path,
		idToPath: mapa,
	}
	cam.idCounter = len(mapa)
	return cam
}

func (r *rpiCam) TakePhoto() (int, error) {
	t := time.Now()
	name := r.photoDir + t.Format("/2006-01-02150405") + strconv.Itoa(r.idCounter+1) + ".jpeg"
	_, err := cmdGetPhoto(name)
	if err != nil {
		return 0, err
	}
	fs, err := os.Stat(name)
	fmt.Println(fs)
	fmt.Println(err)
	if err != nil {
		return 0, err
	}
	fmt.Println(fs)
	r.idCounter++
	r.idToPath[r.idCounter] = fs

	return r.idCounter, nil
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

func cmdGetPhoto(name string) ([]byte, error) {
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
