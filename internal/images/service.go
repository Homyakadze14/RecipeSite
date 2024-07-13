package images

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func Save(path string, image multipart.File) (string, error) {
	dirPath := fmt.Sprintf("./static/%s", path)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp(dirPath, "upload-*.jpg")
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
	if strings.Contains(path, ";") {
		for _, v := range strings.Split(path, ";") {
			err := os.Remove(v)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return os.Remove(path)
}
