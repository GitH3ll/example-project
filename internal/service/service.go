package service

import (
	"context"
	"fmt"
	"github.com/GitH3ll/example-project/internal/config"
	"github.com/GitH3ll/example-project/internal/constants"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"

	"github.com/GitH3ll/example-project/internal/model"
)

type repository interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
	DeleteUser(ctx context.Context, id int64) error
	CheckAuth(ctx context.Context, login, password string) (model.User, error)
}

type imageRepository interface {
	AddImage(ctx context.Context, modelImage model.Image) error
	GetImages(ctx context.Context, userID int) ([]model.Image, error)
}

type fileStorage interface {
	PutObject(ctx context.Context, image model.Image) error
	GetUrls(ctx context.Context, images []model.Image) ([]string, error)
}

type Controller struct {
	repo      repository
	imageRepo imageRepository
	cfg       *config.Config
	minio     fileStorage
}

func NewController(repo repository, imageRepo imageRepository, cfg *config.Config, m fileStorage) *Controller {
	return &Controller{
		repo:      repo,
		imageRepo: imageRepo,
		cfg:       cfg,
		minio:     m,
	}
}

func (c *Controller) AddUser(ctx context.Context, user model.User) error {
	return c.repo.AddUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int64) (model.User, error) {
	user, err := c.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	images, err := c.imageRepo.GetImages(ctx, user.ID)
	if err != nil {
		return model.User{}, err
	}

	urls, err := c.minio.GetUrls(ctx, images)
	if err != nil {
		return model.User{}, err
	}

	user.ImageUrls = urls

	return user, nil
}

func (c *Controller) UpdateUser(ctx context.Context, user model.User) error {
	id := ctx.Value(constants.IdCtxKey)

	if id != user.ID {
		return fmt.Errorf("users do not match")
	}

	return c.repo.UpdateUser(ctx, user)
}

func (c *Controller) DeleteUser(ctx context.Context, id int64) error {
	return c.repo.DeleteUser(ctx, id)
}

func (c *Controller) Authorize(ctx context.Context, login, password string) (string, error) {
	user, err := c.repo.CheckAuth(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authorize user: %w", err)
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		Subject:   "authorized",
		Audience:  []string{"1"},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        strconv.Itoa(int(user.ID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.cfg.JWTKeyword))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (c *Controller) AddFile(ctx context.Context, image model.Image) error {
	imageName, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate image name: %w", err)
	}

	image.Name = imageName.String() + image.Extension

	err = c.imageRepo.AddImage(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to save image data to db: %w", err)
	}

	err = c.minio.PutObject(ctx, image)
	if err != nil {
		return fmt.Errorf("failed to put image to fileStore: %w", err)
	}

	return err
}
