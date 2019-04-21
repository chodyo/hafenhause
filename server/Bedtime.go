package hafenhause

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"time"

	"cloud.google.com/go/firestore"
)

type bedtimedb interface {
	create(map[docpath]interface{}) error
	read(...docpath) ([]*firestore.DocumentSnapshot, error)
	update(map[docpath]interface{}) error
	delete(...docpath) error

	getBedtimes(...docpath) ([]bedtime, error)

	combine(docpath, ...string) docpath
	// docref(docpath, ...string) *firestore.DocumentRef
}

type bedtime struct {
	Name    *string    `json:"name,omitempty"`
	Hour    int        `json:"hour"`
	Minute  int        `json:"minute"`
	Updated *time.Time `json:"updated,omitempty"`
}

var (
	errAlreadyExists         = errors.New("Resource already exists")
	errNotFound              = errors.New("Resource not found")
	errUnsupportedHTTPMethod = errors.New("Unsupported HTTP method")
)

// Bedtime is the web interface between the client and the raw bedtime data
func Bedtime(w http.ResponseWriter, r *http.Request) {
	db := newHafenhausedb()
	processRequest(w, r, db)
}

func processRequest(w http.ResponseWriter, r *http.Request, db bedtimedb) {
	var responseBody []byte
	var err error
	defer func() {
		if err != nil {
			statusCode := statusCodeFromErr(err)
			w.WriteHeader(statusCode)
			panic(err)
		}
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBody)
	}()

	switch r.Method {
	// CREATE
	case http.MethodPost:
		name := path.Base(r.URL.Path)
		switch name {
		case brannigan, malcolm:
			var defaultsList []bedtime
			defaultsList, err = db.getBedtimes(BedtimeDefaultsPath)
			if err != nil {
				return
			}

			defaults := defaultsList[0]

			path := db.combine(BedtimePeoplePath, name)

			now := time.Now()
			bedtime := bedtime{
				Hour:    defaults.Hour,
				Minute:  defaults.Minute,
				Updated: &now,
			}

			err = db.create(map[docpath]interface{}{path: bedtime})

		default:
			err = errNotFound
		}

	// READ
	case http.MethodGet:
		docpaths := []docpath{
			db.combine(BedtimePeoplePath, brannigan),
			db.combine(BedtimePeoplePath, malcolm),
		}
		var bedtimes []bedtime
		if bedtimes, err = db.getBedtimes(docpaths...); err != nil {
			return
		}
		responseBody, err = json.Marshal(bedtimes)

	// UPDATE
	case http.MethodPut:
		var bedtimes []bedtime
		if err = json.NewDecoder(r.Body).Decode(&bedtimes); err != nil {
			return
		}

		var updateEntries map[docpath]interface{}
		for _, bedtime := range bedtimes {
			path := db.combine(BedtimePeoplePath, *bedtime.Name)
			updateEntries[path] = bedtime
		}

		err = db.update(updateEntries)

	// DELETE
	case http.MethodDelete:
		name := path.Base(r.URL.Path)
		switch name {
		case brannigan, malcolm:
			path := db.combine(BedtimePeoplePath, name)
			err = db.delete(path)
		default:
			err = errNotFound
		}

	default:
		err = errUnsupportedHTTPMethod
	}
}

func statusCodeFromErr(e error) int {
	switch e {
	case errAlreadyExists, errUnsupportedHTTPMethod:
		return http.StatusBadRequest

	case errNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError

	}
}

func (h hafenhausedb) getBedtimes(paths ...docpath) ([]bedtime, error) {
	docsnaps, err := h.read(paths...)
	if err != nil {
		return nil, err
	} else if len(docsnaps) != len(paths) {
		return nil, errNotFound
	}

	var bedtimes []bedtime

	for _, ds := range docsnaps {
		var bedtime bedtime

		if err := ds.DataTo(&bedtime); err != nil {
			return nil, err
		}

		bedtime.Name = &ds.Ref.ID

		bedtimes = append(bedtimes, bedtime)
	}

	return bedtimes, nil
}
