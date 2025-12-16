package storage

import "context"

type ObjectStorage interface {
	Upload(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error)
	Update(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error)
	Delete(ctx context.Context, fileName string, bucketName string) error
}
