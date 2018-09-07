package server

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Name   string `json:"name"`
	UserID string `json:"userID"`
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
		UserID: fmt.Sprintf("%x", sha256.Sum224([]byte((name + strconv.FormatInt(time.Now().UnixNano(), 10))))),
	}

	userStore[user.UserID] = user

	return user
}

func RemoveUser(userID string) {

	delete(userStore, userID)
}

func GetUser(userID string) User {
	return userStore[userID]
}
