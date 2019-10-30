package conv

import "time"

const ISO8601 = "2006-01-02T15:04:05-0700"

func ParseISO8601(in string) (time.Time, error) {
	return time.Parse(time.RFC3339, in)
}

func FormatISO8601(in time.Time) string {
	return in.Format(time.RFC3339)
}
