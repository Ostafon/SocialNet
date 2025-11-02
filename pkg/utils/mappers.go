package utils

import (
	"fmt"
	"strconv"
)

func StringToUint(id string) (uint, error) {
	var UId, err = strconv.Atoi(id)
	if err != nil {
		return 0, ErrorHandler(err, "cannot convert string")
	}

	return uint(UId), nil
}

func UintToString(id uint) string {

	return fmt.Sprint(id)
}
