package userStore

type Authorization interface {
	CreateUser(username, password string, passcode string) error
	CheckUser(username, password string) bool
	CheckUserExist(username string) bool
}
