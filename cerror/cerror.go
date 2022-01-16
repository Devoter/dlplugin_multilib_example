package cerror

import "fmt"

type CError int

func (cErr CError) Error() string {
	return fmt.Sprintf("C error code is: %d", cErr)
}

func WrapCError(code int) error {
	if code == 0 {
		return nil
	}

	return CError(code)
}
