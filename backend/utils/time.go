package utils

import "time"

func UnixToIso(unix int64) string {
	t := time.Unix(unix, 0).UTC()
	s := t.Format(time.RFC3339Nano)

	return s
}
