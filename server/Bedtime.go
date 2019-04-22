package hafenhause

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chodyo/hafenhause/nosqldb"

	"github.com/codemodus/parth/v2"
)

type bedtimedbContract interface {
	createDefaultBedtime(string) error
	getBedtimes(string) ([]bedtime, error)
	updateBedtime(string, bedtime) error
	deleteBedtime(string) error
}

type bedtime struct {
	Name    *string    `json:"name,omitempty"`
	Hour    int        `json:"hour"`
	Minute  int        `json:"minute"`
	Updated *time.Time `json:"updated,omitempty"`
}

var (
	errBadRequest = errors.New("Bad request")
)

// Bedtime is the web interface between the client and the raw bedtime data
func Bedtime(w http.ResponseWriter, r *http.Request) {
	db := newBedtimedb()

	var responseBody []byte
	var err error

	defer func() {
		if err != nil {
			http.Error(w, err.Error(), statusCodeFromErr(err))
			return
		}

		if responseBody != nil {
			fmt.Fprint(w, string(responseBody))
		}
	}()

	log.Printf("Processing request with URL=%s and body=%+v\n", r.URL.String(), r.Body)

	responseBody, err = processRequest(w, r, db)
}

func processRequest(w http.ResponseWriter, r *http.Request, db bedtimedbContract) (responseBody []byte, err error) {
	switch r.Method {
	// CREATE
	case http.MethodPost:
		var name string
		name, err = getNamePathParam(r.URL)

		err = db.createDefaultBedtime(name)

	// READ
	case http.MethodGet:
		var name string
		name = getNameQueryParam(r.URL)

		var bedtimes []bedtime
		bedtimes, err = db.getBedtimes(name)
		if err != nil {
			return
		}

		responseBody, err = json.Marshal(bedtimes)

	// UPDATE
	case http.MethodPut:
		var name string
		name, err = getNamePathParam(r.URL)

		var bedtime bedtime
		err = json.NewDecoder(r.Body).Decode(&bedtime)
		if err != nil {
			err = errBadRequest
			return
		}

		err = db.updateBedtime(name, bedtime)

	// DELETE
	case http.MethodDelete:
		var name string
		name, err = getNamePathParam(r.URL)

		err = db.deleteBedtime(name)

	default:
		err = http.ErrNotSupported
	}

	return
}

func getNamePathParam(URL *url.URL) (name string, err error) {
	path := URL.Path

	if err = parth.Segment(path, 0, &name); err != nil {
		err = errBadRequest
	}

	log.Printf("Found path param name=%s\n", name)

	return
}

func getNameQueryParam(URL *url.URL) (name string) {
	queryValues := URL.Query()

	if name = queryValues.Get("name"); name == "" {
		name = "*"
	}

	log.Printf("Found query param name=%s\n", name)

	return
}

func statusCodeFromErr(e error) int {
	switch e {
	case nosqldb.ErrAlreadyExists, http.ErrNotSupported:
		return http.StatusBadRequest

	case nosqldb.ErrNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
