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
	Subject   string     `json:"subject"`
	Date      *time.Time `json:"date,omitempty"`
	Score     int        `json:"score,omitempty"`
	CarryOver bool       `json:"carryOver,omitempty"`

	db *firestore.CollectionRef
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
	report.db = client.Collection(collection)
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

	var bedtime time.Time

	// explicit dates override everything
	if r.Date != nil {
		bedtime = *r.Date
	} else {
		bedtime = r.GetLastOrDefault()

		if r.Score != 0 {
			bedtime = bedtime.Add(time.Duration(r.Score) * time.Minute)
		}

		if bedtime.Before(time.Now()) {
			bedtime = bedtime.Add(24 * time.Hour)
		}

		defaultBedtime := GetDefaultBedtime()
		if bedtime.After(defaultBedtime) {
			bedtime = defaultBedtime
		}
	}

	_, err := r.db.Doc(r.Subject).Set(ctx, map[string]time.Time{
		"bedtime": bedtime,
	})

	return err
}

// GetLastOrDefault queries the database and returns a default value in the
// future if the entry does not exist
func (r *Report) GetLastOrDefault() time.Time {
	ctx := context.Background()

	subjectDoc, err := r.db.Doc(r.Subject).Get(ctx)
	if !subjectDoc.Exists() || err != nil {
		return GetDefaultBedtime()
	}

	var lastBedtime map[string]time.Time
	subjectDoc.DataTo(&lastBedtime)

	return lastBedtime["bedtime"]
}

// GetAllBedtimes returns the times and dates of all stored bedtimes for
// relevant subjects, i.e. Bran and Mal
func GetAllBedtimes() ([]Report, error) {
	ctx := context.Background()

	docsnaps, err := client.GetAll(ctx, []*firestore.DocumentRef{
		client.Doc("hafenhaus/" + brannigan),
		client.Doc("hafenhaus/" + malcolm),
	})
	if err != nil {
		return nil, err
	}

	var bedtimes []Report

	for _, ds := range docsnaps {
		var report Report

		report.Subject = ds.Ref.ID

		bedtimeEntity, err := ds.DataAt("bedtime")
		if err != nil {
			return []Report{}, err
		}
		bedtime, ok := bedtimeEntity.(time.Time)
		if !ok {
			return []Report{}, err
		}
		report.Date = &bedtime

		bedtimes = append(bedtimes, report)
	}

	return bedtimes, nil
}

// GetDefaultBedtime will return the default (maximum) next bedtime. If today's
// bedtime is in the past, it will add 1 day to ensure it is referencing
// tomorrow's bedtime
func GetDefaultBedtime() time.Time {
	loc, err := time.LoadLocation(defaultTimeZone)

	if err != nil {
		log.Printf("ERROR loading time location! %v", err)
	}

	now := time.Now().In(loc)

	bedtime := time.Date(now.Year(), now.Month(), now.Day(), defaultBedtimeHour,
		defaultBedtimeMinute, int(0), int(0), loc)

	if bedtime.Before(now) {
		bedtime = bedtime.Add(24 * time.Hour)
	}

	return bedtime
}
