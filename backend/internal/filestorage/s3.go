package filestorage

import (
	"fmt"
	"io"
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

func (s *S3Storage) saveToS3(urlCh chan<- string, errCh chan<- error, photo io.ReadSeeker, contentType string) {
	uid := uuid.New().String()
	_, err := s.PutObject(&s3.PutObjectInput{
		Body:        photo,
		Bucket:      s.bucket,
		Key:         aws.String(uid),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		errCh <- err
	}

	urlCh <- fmt.Sprintf("%s/%s/%s", s.Endpoint, *s.bucket, uid) + ";"
}

func (s *S3Storage) Save(photos []io.ReadSeeker, contentType string) (string, error) {
	urls := ""

	urlChan := make(chan string)
	errChan := make(chan error)

	defer close(urlChan)
	defer close(errChan)

	for _, photo := range photos {
		go s.saveToS3(urlChan, errChan, photo, contentType)
	}

	for i := 0; i < len(photos); i++ {
		select {
		case url := <-urlChan:
			urls += url
		case err := <-errChan:
			return "", err
		}
	}

	return urls, nil
}

func getFileName(path string) string {
	elems := strings.Split(path, "/")
	return elems[len(elems)-1]
}

func (s *S3Storage) Remove(path string) error {
	if path == s.defaultIconUrl {
		return nil
	}

	urls := strings.Split(path, ";")

	for _, url := range urls {
		if url == "" {
			continue
		}

		_, err := s.DeleteObject(&s3.DeleteObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(getFileName(url)),
		})
		if err != nil {
			return fmt.Errorf("S3 - Remove - s.DeleteObject: %w", err)
		}
	}

	return nil
}
