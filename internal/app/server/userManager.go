package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
)

func newUserStore(userFile, passcode string) userStore {
	users := map[string][]byte{}
	jsonFile, err := os.OpenFile(userFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Open userFile error")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &users)
	return userStore{
		userFile:     userFile,
		userPassword: users,
		passcode:     passcode,
	}
}

type userStore struct {
	userFile     string
	userPassword map[string][]byte
	passcode     string
}

func (u userStore) checkUserExist(username string) bool {
	_, ok := u.userPassword[username]
	return ok
}

func (u *userStore) CreateUser(username, password string, passcode string) error {
	fmt.Println(passcode, u.passcode)
	if u.passcode != passcode {
		return errors.New("Wrong passcode")
	}
	if u.checkUserExist(username) {
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

func (u *userStore) CheckUser(username, password string) bool {
	if !u.checkUserExist(username) {
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
func (u userStore) save() error {
	jsonFile, err := os.OpenFile(u.userFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, _ := json.Marshal(&u.userPassword)
	_, err = jsonFile.Write(bytes)
	return nil
}
