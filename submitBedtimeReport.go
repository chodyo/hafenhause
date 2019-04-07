package hafenhaus

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/chodyo/hafenhaus/bedtime"
)

const (
	decodeMessage   = "Could not decode request: %v"
	validateMessage = "Could not parse request: %v"
	persistMessage  = "Could not persist data: %v"
	successMessage  = "Success! Report saved: %+v"
)

// SubmitBedtimeReport receives a bedtimeReport and saves it to a FireStore
func SubmitBedtimeReport(w http.ResponseWriter, r *http.Request) {
	report, err := bedtime.NewReportFromRequest(r)
	if err != nil {
		log.Printf(decodeMessage, err)
		writeResponse(w, http.StatusBadRequest, decodeMessage, err)
		return
	}

	if errs := report.Validate(); len(errs) > 0 {
		log.Printf(validateMessage, errs)
		writeResponse(w, http.StatusBadRequest, validateMessage, errs)
		return
	}

	if err := report.Save(); err != nil {
		log.Printf(persistMessage, err)
		writeResponse(w, http.StatusInternalServerError, persistMessage, err)
		return
	}

	log.Printf(successMessage, report)
	writeResponse(w, http.StatusOK, successMessage, report)
}

func writeResponse(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, fmt.Sprintf(msg, args...))
}
