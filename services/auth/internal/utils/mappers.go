package utils

import "strconv"

func StringToUint(id string) (uint, error) {
	var UId, err = strconv.Atoi(id)
	if err != nil {
		return 0, ErrorHandler(err, "cannot convert string")
	}

	return uint(UId), nil
}
