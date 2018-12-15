package cb_errors

import (
	"fmt"
)

var (
	CONTACT_NOT_FOUND error = fmt.Errorf("Requested contact not found")
)
