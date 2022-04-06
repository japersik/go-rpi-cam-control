package jsonStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
)

func NewUserStore(config *Config) *UserStore {
	users := map[string][]byte{}
	jsonFile, err := os.OpenFile(config.AuthFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Open userFile error")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &users)
	return &UserStore{
		userFile:     config.AuthFile,
		userPassword: users,
		passcode:     config.NewUserKey,
	}
}

type UserStore struct {
	userFile     string
	userPassword map[string][]byte
	passcode     string
}

func (u UserStore) CheckUserExist(username string) bool {
	_, ok := u.userPassword[username]
	return ok
}

func (u *UserStore) CreateUser(username, password string, passcode string) error {
	fmt.Println(passcode, u.passcode)
	if u.passcode != passcode {
		return errors.New("Wrong passcode")
	}
	if u.CheckUserExist(username) {
		return errors.New("User exists")
	}
	if len(password) < 6 {
		return errors.New("Password should be at least 6 letters long")
	}
	hPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.userPassword[username] = hPass
	u.save()
	return nil
}

func (u *UserStore) CheckUser(username, password string) bool {
	if !u.CheckUserExist(username) {
		return false
	}
	chash := u.userPassword[username]
	err := bcrypt.CompareHashAndPassword(chash, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func getHashingCost(hashedPassword []byte) int {
	cost, _ := bcrypt.Cost(hashedPassword) // Игнорировать обработку ошибок для простоты
	return cost
}

func (u UserStore) save() error {
	jsonFile, err := os.OpenFile(u.userFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, _ := json.Marshal(&u.userPassword)
	_, err = jsonFile.Write(bytes)
	return nil
}
