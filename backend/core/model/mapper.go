package model

import (
	"fmt"
)

func TruncateID(ID string) string {
	if IsUUID(ID) {
		return fmt.Sprintf("%.8s", ID)
	}

	if len(ID)/3 >= 8 {
		return fmt.Sprintf("%.8s", ID)
	}

	return fmt.Sprintf("%.6s", ID)
}
