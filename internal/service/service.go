package service

import "github.com/GitH3ll/example-project/internal/model"

type repository interface {
	AddUser(user model.User) error
	GetUser(id int) (model.User, error)
}

type User struct {
	repo repository
}

func NewUserService(repo repository) *User {
	return &User{repo: repo}
}

func (u *User) AddUser(user model.User) error {
	return u.repo.AddUser(user)
}

func (u *User) GetUser(id int) (model.User, error) {
	return u.repo.GetUser(id)
}
