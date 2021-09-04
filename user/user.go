package user

import (
	"time"
)

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Token string `json:"token"`
	DBMaster string `json:"db_master"`
	ID int `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	}

func New(name string, age int, token string, dbMaster string) *User {
	return &User{
		Name: name,
		Age: age,
		Token: token,
		DBMaster: dbMaster,
		}
}
