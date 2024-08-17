package common

import "errors"

var (
	ErrServerError      = errors.New("server error")
	ErrUrlParam         = errors.New("ID must be provided")
	ErrRecipeIDType     = errors.New("ID must be integer")
	ErrNoPermissions    = errors.New("no permissions")
	ErrCacheKeyNotFound = errors.New("key not found")
	ErrLoginProvided    = errors.New("login must be provided in url")
	ErrHudgeFiles       = errors.New("files are too hudge")
	ErrImageType        = errors.New("files must be images")
)

func HavePermisson(needID, haveID int) bool {
	return needID == haveID
}
