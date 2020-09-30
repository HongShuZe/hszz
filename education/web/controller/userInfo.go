package controller

import "github.com/hszz/education/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	IsAdmin   string
}

var users []User

func init() {
	admin := User{LoginName: "admin", Password: "123456", IsAdmin: "T"}
	king := User{LoginName: "king", Password: "123456", IsAdmin: "T"}
	queen := User{LoginName: "queen", Password: "123456", IsAdmin: "F"}
	joker := User{LoginName: "joker", Password: "123456", IsAdmin: "F"}

	users = append(users, admin)
	users = append(users, king)
	users = append(users, queen)
	users = append(users, joker)
}
