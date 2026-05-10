package minio

import (
"context"
"fmt"
"time"

"github.com/minio/minio-go/v7"
"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "listings"

type Storage struct {
client *minio.Client
}

func NewStorage(endpoint, accessKey, secretKey string) (*Storage, error) {
client, err := minio.New(endpoint, &minio.Options{
Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
Secure: false,
})
if err != nil {
return nil, fmt.Errorf("создание MinIO клиента: %w", err)
}

return &Storage{client: client}, nil
}

func (s *Storage) EnsureBucket(ctx context.Context) error {
exists, err := s.client.BucketExists(ctx, bucketName)
if err != nil {
return fmt.Errorf("проверка бакета: %w", err)
}
if exists {
return nil
}

if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
return fmt.Errorf("создание бакета: %w", err)
}

policy := fmt.Sprintf(`{
"Version":"2012-10-17",
"Statement":[{
"Effect":"Allow",
"Principal":{"AWS":["*"]},
"Action":["s3:GetObject"],
"Resource":["arn:aws:s3:::%s/*"]
}]
}`, bucketName)

if err := s.client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
return fmt.Errorf("установка политики бакета: %w", err)
}

return nil
}

func (s *Storage) GetUploadURL(ctx context.Context, listingID, filename, contentType string) (uploadURL, fileURL string, err error) {
objectName := fmt.Sprintf("%s/%s", listingID, filename)

presignedURL, err := s.client.PresignedPutObject(ctx, bucketName, objectName, 15*time.Minute)
if err != nil {
return "", "", fmt.Errorf("генерация presigned URL: %w", err)
}

fileURL = fmt.Sprintf("http://%s/%s/%s", s.client.EndpointURL().Host, bucketName, objectName)

return presignedURL.String(), fileURL, nil
}
