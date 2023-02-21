package service

import (
	"context"
	"github.com/GitH3ll/example-project/internal/model"
)

type repository interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
}

type Controller struct {
	repo repository
}

func NewController(repo repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) AddUser(ctx context.Context, user model.User) error {
	return c.repo.AddUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int64) (model.User, error) {
	return c.repo.GetUser(ctx, id)
}

func (c *Controller) UpdateUser(ctx context.Context, user model.User) error {
	return c.repo.UpdateUser(ctx, user)
}
