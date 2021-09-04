package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type UserRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *UserRepo {
	return &UserRepo{conn}
}

func (r *UserRepo) WithID(id int) (*User, error) {
	user := User{}
	err := r.conn.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, ErrNoRecord
	}

	return &user, err
}

func (r *UserRepo) Store(user *User) error {
	return r.conn.Create(user).Error
}
