package util

import "time"

// Now returns time with microseconds
func Now() time.Time {
	return time.Unix(0, time.Now().UnixNano()/1e6*1e6)
}
