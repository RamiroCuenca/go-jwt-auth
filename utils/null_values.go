package utils

import (
	"time"

	"github.com/lib/pq"
)

// Convertes an empty date into a null value to be stored on the database.
func DateToNull(d time.Time) pq.NullTime {
	t := pq.NullTime{Time: d}

	// If time is not zero, return that it's valid
	if !t.Time.IsZero() {
		t.Valid = true
	}

	return t
}
