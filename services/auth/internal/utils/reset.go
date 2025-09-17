package utils

import "github.com/google/uuid"

func GenerateUUID() (string, error) {
	id, err := uuid.NewUUID()

	return id.String(), err
}
