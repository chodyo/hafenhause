package bedtimes

import (
	"time"
)

// Report summarizes the bedtime routine for a single person on a single night
type Report struct {
	Subject person     `json:"subject,omitempty"`
	Date    *time.Time `json:"date,omitempty"`
	Score   int        `json:"score,omitempty"`
	Notes   string     `json:"notes,omitempty"`
}
