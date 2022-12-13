package model

import (
	"fmt"

	"github.com/google/uuid"
)

func TruncateID(ID string) string {
	if isUUID(ID) {
		return fmt.Sprintf("%.8s", ID)
	}

	if len(ID)/3 >= 8 {
		return fmt.Sprintf("%.8s", ID)
	}

	return fmt.Sprintf("%.6s", ID)
}

func isUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}

	return true
}
