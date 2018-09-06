package server

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

var userStore map[string]User

func InitUserStore(force bool) {
	if userStore == nil || force {
		userStore = make(map[string]User)
	}
}

func CreateUser(name string) User {

	InitUserStore(false)

	user := User{
		Name:   name,
		UserId: fmt.Sprintf("%x", sha256.Sum224([]byte((name + strconv.FormatInt(time.Now().UnixNano(), 10))))),
	}

	userStore[user.UserId] = user

	return user
}

func RemoveUser(userId string) {

	delete(userStore, userId)
}

func GetUser(userId string) User {
	return userStore[userId]
}
