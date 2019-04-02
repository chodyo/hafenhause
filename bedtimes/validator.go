package bedtimes

import (
	"fmt"
	"time"
)

// ValidateReport makes sure the Report object has sane defaults
func ValidateReport(report *Report) []error {
	var errs []error

	// make sure we know who the report is for
	if report.Subject < cody || report.Subject > malcolm {
		errs = append(errs, fmt.Errorf("invalid person"))
	}

	reportDate := report.Date

	// assume report was for today if date not specified
	if reportDate == nil {
		now := time.Now()
		reportDate = &now
	}

	// round date to day
	roundedDate := reportDate.Truncate(24 * time.Hour).UTC()
	report.Date = &roundedDate

	return errs
}
