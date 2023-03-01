package filestore

import (
	"context"
	"github.com/GitH3ll/example-project/internal/model"
	"github.com/minio/minio-go/v7"
	"time"
)

type Minio struct {
	minio  *minio.Client
	bucket string
}

func NewMinio(minioClient *minio.Client, bucket string) *Minio {
	return &Minio{
		minio:  minioClient,
		bucket: bucket,
	}
}

func (m *Minio) PutObject(ctx context.Context, image model.Image) error {
	_, err := m.minio.PutObject(ctx, m.bucket, image.Name, image.Data, -1, minio.PutObjectOptions{})

	return err
}

func (m *Minio) GetUrls(ctx context.Context, images []model.Image) ([]string, error) {
	var urls []string

	for i := range images {
		url, err := m.minio.PresignedGetObject(ctx, m.bucket, images[i].Name, time.Hour*24, nil)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url.String())
	}

	return urls, nil
}
