package filestorage

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Homyakadze14/RecipeSite/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	*s3.S3
	bucket         *string
	defaultIconUrl string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.S3.ACCESS_KEY, cfg.S3.SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(cfg.S3.ENDPOINT),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(newSession)

	return &S3Storage{s3Client, &cfg.S3.BUCKET_NAME, cfg.S3.DEFAULT_ICON_URL}, nil
}

func (s *S3Storage) Save(image multipart.File, contentType string) (string, error) {
	uid := uuid.New().String()
	_, err := s.PutObject(&s3.PutObjectInput{
		Body:        image,
		Bucket:      s.bucket,
		Key:         aws.String(uid),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.Endpoint, *s.bucket, uid), nil
}

func getFileName(path string) string {
	elems := strings.Split(path, "/")
	return strings.Split(path, "/")[len(elems)-1]
}

func (s *S3Storage) Remove(path string) error {
	if path == s.defaultIconUrl {
		return nil
	}

	if strings.Contains(path, ";") {
		for _, v := range strings.Split(path, ";") {
			if v == "" {
				continue
			}

			_, err := s.DeleteObject(&s3.DeleteObjectInput{
				Bucket: s.bucket,
				Key:    aws.String(getFileName(v)),
			})
			if err != nil {
				return err
			}
		}
		return nil
	}

	_, err := s.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(getFileName(path)),
	})

	return err
}
