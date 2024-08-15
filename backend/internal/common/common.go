package common

import "errors"

var (
	ErrServerError  = errors.New("server error")
	ErrUrlParam     = errors.New("ID must be provided")
	ErrRecipeIDType = errors.New("ID must be integer")
)

func HavePermisson(needID, haveID int) bool {
	return needID != haveID
}
