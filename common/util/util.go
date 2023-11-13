package util

import "github.com/google/uuid"

func IsValidUUIDv4(id string) bool {
	u, err := uuid.Parse(id)
	return err == nil && u.Version().String() == "VERSION_4"
}
