package hafenhause

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type bedtimedb interface {
	createBedtimes([]string) error
	readBedtimes([]string) ([]bedtime, error)
	updateBedtimes([]bedtime) error
	deleteBedtimes([]string) error
}

type bedtime struct {
	Name   string `json:"name"`
	Hour   int    `json:"hour"`
	Minute int    `json:"minute"`
}

// Bedtime is the web interface between the client and the raw bedtime data
// GET: returns a list of
// PUT: receives a list of bedtime values and updates
func Bedtime(w http.ResponseWriter, r *http.Request) {
	db := newHafenhausedb()
	processRequest(w, r, db)
}

func processRequest(w http.ResponseWriter, r *http.Request, db bedtimedb) {
	var responseBody []byte
	var err error
	defer func() {
		if err != nil {
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
		err = db.createBedtimes([]string{name})

	// READ
	case http.MethodGet:
		names := []string{brannigan, malcolm}
		var bedtimes []bedtime
		if bedtimes, err = db.readBedtimes(names); err != nil {
			return
		}
		responseBody, err = json.Marshal(bedtimes)

	// UPDATE
	case http.MethodPut:
		var bedtimes []bedtime
		if err = json.NewDecoder(r.Body).Decode(&bedtimes); err != nil {
			return
		}
		err = db.updateBedtimes(bedtimes)

	// DELETE
	case http.MethodDelete:
		name := path.Base(r.URL.Path)
		err = db.deleteBedtimes([]string{name})

	default:
		err = fmt.Errorf("Unsupported HTTP method")
	}
}
