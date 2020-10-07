package quinyx

import (
	"strconv"
	"time"
)

// Timestamp type
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON into a Timestamp
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		t.Time = time.Unix(i, 0)
	} else {
		t.Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
	}
	return
}

// Equal compares timestamps
func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}
