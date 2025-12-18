package usecases

import "github.com/google/uuid"

func ptrBool(value bool) *bool {
	return &value
}

func ptrStr(value string) *string {
	return &value
}

func ptrInt(value int) *int {
	return &value
}

func ptrUUID(value uuid.UUID) *uuid.UUID {
	return &value
}
