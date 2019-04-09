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
		http.Error(w, decodeMessage, http.StatusBadRequest)
		return
	}

	if errs := report.Validate(); len(errs) > 0 {
		log.Printf(validateMessage, errs)
		http.Error(w, validateMessage, http.StatusBadRequest)
		return
	}

	if err := report.Save(); err != nil {
		log.Printf(persistMessage, err)
		http.Error(w, persistMessage, http.StatusInternalServerError)
		return
	}

	log.Printf(successMessage, report)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, fmt.Sprintf(successMessage, report))
}
