package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ramdhanrizkij/next-store-api/internal/config"
)

type StorageService interface {
	UploadFile(ctx context.Context, folder string, filename string, reader io.Reader, size int64, contentType string) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
	GetFileURL(objectName string) string
}

type minioStorageService struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewMinioStorageService(cfg config.MinioConfig) (StorageService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ensure bucket exists
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.Printf("[WARN] Failed to check if minio bucket exists: %v. Will continue anyway.", err)
	} else if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("[WARN] Failed to create minio bucket '%s': %v", cfg.Bucket, err)
		} else {
			log.Printf("[INFO] Created minio bucket '%s'", cfg.Bucket)
			// Set public read policy for assets
			policy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::%s/*"]
					}
				]
			}`, cfg.Bucket)
			_ = client.SetBucketPolicy(ctx, cfg.Bucket, policy)
		}
	}

	publicURL := strings.TrimRight(cfg.PublicURL, "/")
	if publicURL == "" {
		scheme := "http"
		if cfg.UseSSL {
			scheme = "https"
		}
		publicURL = fmt.Sprintf("%s://%s/%s", scheme, cfg.Endpoint, cfg.Bucket)
	}

	return &minioStorageService{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: publicURL,
	}, nil
}

func (s *minioStorageService) UploadFile(
	ctx context.Context,
	folder string,
	filename string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	ext := filepath.Ext(filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.NewString(), ext)

	objectPath := uniqueName
	if folder != "" {
		folder = strings.Trim(folder, "/")
		objectPath = fmt.Sprintf("%s/%s", folder, uniqueName)
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := s.client.PutObject(ctx, s.bucket, objectPath, reader, size, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upload object to minio: %w", err)
	}

	return s.GetFileURL(objectPath), nil
}

func (s *minioStorageService) DeleteFile(ctx context.Context, objectName string) error {
	// Extract object path if full URL is passed
	if strings.Contains(objectName, s.bucket) {
		idx := strings.Index(objectName, s.bucket)
		objectName = strings.TrimPrefix(objectName[idx+len(s.bucket):], "/")
	}

	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *minioStorageService) GetFileURL(objectName string) string {
	objectName = strings.TrimPrefix(objectName, "/")
	return fmt.Sprintf("%s/%s", s.publicURL, objectName)
}
