package common

import "errors"

var (
	ErrServerError = errors.New("server error")
)

func HavePermisson(needID, haveID int) bool {
	return needID != haveID
}
