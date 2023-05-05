package common

import (
	"time"
)

// time format template(UnmarshalJSON and MarshalJSON)
const (
	TimeFormat        = "2006-01-02"
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

// JSONTime is time
type JSONTime time.Time

// UnmarshalJSON for JSON Time
func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = JSONTime(now)
	return
}

// MarshalJSON for JSON Time
func (t JSONTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

// String for JSON Time
func (t JSONTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
