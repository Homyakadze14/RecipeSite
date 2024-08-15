package common

func HavePermisson(needID, haveID int) bool {
	return needID != haveID
}
