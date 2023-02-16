package service

import "github.com/GitH3ll/example-project/internal/model"

type repository interface {
	AddUser(user model.User) error
	GetUser(id int) (model.User, error)
}

type Controller struct {
	repo repository
}

func NewController(repo repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) AddUser(user model.User) error {
	return c.repo.AddUser(user)
}

func (c *Controller) GetUser(id int) (model.User, error) {
	return c.repo.GetUser(id)
}
