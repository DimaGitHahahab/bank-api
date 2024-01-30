package model

import "time"

type UserInfo struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}
