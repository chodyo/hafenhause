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
	createDefaultBedtime(name string) (err error)
	getBedtimes(name string) (bedtimes []bedtime, err error)
	updateBedtime(name string, bedtime bedtime) (err error)
	deleteBedtime(name string) (err error)
}

type bedtime struct {
	Name    *string   `json:"name" firestore:"name,omitempty"`
	Hour    int       `json:"hour" firestore:"hour"`
	Minute  int       `json:"minute" firestore:"minute"`
	Updated time.Time `json:"updated" firestore:"updated"`
}

const accessWhitelist = "https://hafenhause.appspot.com"

var errBadRequest = errors.New("Bad request")

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

	// The basic CORS origin header needs to be on every response
	w.Header().Set("Access-Control-Allow-Origin", accessWhitelist)

	log.Printf("Processing request with URL=%s and body=%+v\n", r.URL.String(), r.Body)

	responseBody, err = processRequest(&w, r, db)
}

func processRequest(w *http.ResponseWriter, r *http.Request, db bedtimedbContract) (responseBody []byte, err error) {
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

	// OPTIONS
	case http.MethodOptions:
		// These advanced options tell browsers how to handle CORS. They do not
		// need to be set on every request.
		(*w).Header().Set("Vary", "Origin")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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
