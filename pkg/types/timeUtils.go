package types

import (
	"fmt"
	"time"
)

func ExpiryString(expires time.Time) string {
	return fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", expires.Unix(), expires.String())
}
