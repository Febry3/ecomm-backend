package storage

import "context"

type ObjectStorage interface {
	Upload(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error)
	Delete(ctx context.Context, fileName string) error
}
