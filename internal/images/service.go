package images

import (
	"io"
	"mime/multipart"
	"os"
)

func Save(image multipart.File) (string, error) {
	tempFile, err := os.CreateTemp("./static", "upload-*.jpg")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(image)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)

	return tempFile.Name(), nil
}

func Remove(path string) error {
	return os.Remove(path)
}
