package images

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Homyakadze14/RecipeSite/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	*s3.S3
	bucket *string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.S3_ACCESS_KEY, cfg.S3_SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(cfg.S3_ENDPOINT),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(newSession)

	return &S3Storage{s3Client, &cfg.S3_BUCKET_NAME}, nil
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
	if path == config.DefaultIconURL {
		return nil
	}

	if strings.Contains(path, ";") {
		for _, v := range strings.Split(path, ";") {
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

// func Save(path string, image multipart.File) (string, error) {
// 	dirPath := fmt.Sprintf("./static/%s", path)

// 	err := os.MkdirAll(dirPath, os.ModePerm)
// 	if err != nil {
// 		return "", err
// 	}

// 	tempFile, err := os.CreateTemp(dirPath, "upload-*.jpg")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer tempFile.Close()

// 	fileBytes, err := io.ReadAll(image)
// 	if err != nil {
// 		return "", err
// 	}
// 	tempFile.Write(fileBytes)

// 	return tempFile.Name(), nil
// }

// func Remove(path string) error {
// 	if strings.Contains(path, ";") {
// 		for _, v := range strings.Split(path, ";") {
// 			err := os.Remove(v)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}

// 	return os.Remove(path)
// }
