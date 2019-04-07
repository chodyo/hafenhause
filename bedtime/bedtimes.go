package bedtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	collection string = "hafenhaus"

	defaultTimeZone      string = "America/New_York"
	defaultBedtimeHour   int    = 7 + 12
	defaultBedtimeMinute int    = 30

	cody      string = "cody"
	julia     string = "julia"
	brannigan string = "brannigan"
	malcolm   string = "malcolm"
)

var projectID string
var client *firestore.Client

func init() {
	ctx := context.Background()

	projectID = os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("Set Firebase project ID via GCLOUD_PROJECT env variable")
	}

	var err error
	client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}
}

// Report summarizes the bedtime routine for a single person on a single night.
// Specify either the Date to set an absolute bedtime for the next night or
// specify the Score to set a relative bedtime compared to the night before.
// Use the CarryOver field to set the bedtime to be the same the next night as
// it was the night before.
type Report struct {
	Subject   string     `json:"subject,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
	Score     int        `json:"score,omitempty"`
	CarryOver bool       `json:"carryOver,omitempty"`

	db *firestore.DocumentRef
}

// Reporter contains all functions that can be performed on a Report
type Reporter interface {
	Validate() []error
	Save() error
}

// NewReportFromRequest creates a bedtime Report from an http Request
func NewReportFromRequest(r *http.Request) (*Report, error) {
	var report Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		return nil, err
	}
	report.db = client.Collection(collection).Doc(report.Subject)
	return &report, nil
}

// Validate makes sure the Report object has sane defaults
func (r *Report) Validate() []error {
	var errs []error

	// make sure we know who the report is for
	switch strings.ToLower(r.Subject) {
	case cody, julia, brannigan, malcolm:
	default:
		errs = append(errs, fmt.Errorf("Invalid person"))
	}

	// ensure either a Date or Score was specified in the report
	// Date is absolute, Score is relative to the last value
	// CarryOver should be set if Score should be zero
	if r.Date == nil && r.Score == 0 && !r.CarryOver {
		errs = append(errs, fmt.Errorf("Must set date, score, or carryOver"))
	}

	return errs
}

// Save persists data in FireBase
func (r *Report) Save() error {
	ctx := context.Background()

	bedtime, err := r.GetLastBedtime()
	if err != nil {
		return err
	}

	if r.Date != nil {
		bedtime = *r.Date
	}

	bedtime = ValidOrDefault(bedtime)

	if r.Score != 0 {
		bedtime = bedtime.Add(time.Duration(r.Score) * time.Minute).Add(24 * time.Hour)
	}

	_, err = r.db.Set(ctx, map[string]time.Time{
		"bedtime": bedtime,
	})

	return err
}

// GetLastBedtime queries the database and returns a default value in the future
// if the entry does not exist
func (r *Report) GetLastBedtime() (time.Time, error) {
	ctx := context.Background()

	subjectDoc, err := r.db.Get(ctx)
	if !subjectDoc.Exists() {
		return time.Time{}, nil
	} else if err != nil {
		return time.Time{}, fmt.Errorf("Could not access previous bedtime: %+v", err)
	}

	var lastBedtime map[string]time.Time
	subjectDoc.DataTo(&lastBedtime)

	return lastBedtime["bedtime"], nil
}

// ValidOrDefault checks if a given bedtime is within 24 hours or else it gets
// the next future default bedtime
func ValidOrDefault(t time.Time) time.Time {
	loc, err := time.LoadLocation(defaultTimeZone)
	if err != nil {
		log.Printf("ERROR loading time location! %v", err)
	}
	now := time.Now().In(loc)

	if t.Before(now.Add(24*time.Hour)) && t.After(now) {
		return t
	}

	bedtime := time.Date(now.Year(), now.Month(), now.Day(), defaultBedtimeHour,
		defaultBedtimeMinute, int(0), int(0), loc)
	if bedtime.Before(now) {
		bedtime = bedtime.Add(24 * time.Hour)
	}
	return bedtime
}
