package server

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func newUserStore(userFile string) userStore {
	return userStore{
		userFile:     userFile,
		userPassword: map[string][]byte{},
	}
}

type userStore struct {
	userFile     string
	userPassword map[string][]byte
}

func (u userStore) checkUserExist(username string) bool {
	_, ok := u.userPassword[username]
	return ok
}

func (u *userStore) CreateUser(username, password string) error {
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

//func PassWordHashingHandler(w http.ResponseWriter, r *http.Request) {
//	password := "secret"
//	hash, _ := HashPassword(password) // Для простоты игнорировать обработку ошибок
//
//	fmt.Fprintln(w, "Password:", password)
//	fmt.Fprintln(w, "Hash:    ", hash)
//
//	match := CheckPasswordHash(password, hash)
//	fmt.Fprintln(w, "Match:   ", match)
//
//	cost := GetHashingCost([]byte(hash))
//	fmt.Fprintln(w, "Cost:    ", cost)
//
//}
