package hafenhaus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/chodyo/hafenhaus/bedtimes"
)

// SubmitBedtimeReport receives a bedtimeReport and saves it to a FireStore
func SubmitBedtimeReport(w http.ResponseWriter, r *http.Request) {
	var report bedtimes.Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Printf("Could not decode request: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, fmt.Sprintf("Could not decode request: %v", err))
		return
	}

	if errs := bedtimes.ValidateReport(&report); len(errs) > 0 {
		log.Printf("Could not parse request: %+v", errs)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, fmt.Sprintf("Could not parse request: %v", errs))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: check if report day already has an entry in the db

	err := bedtimes.Save(&report)
	if err != nil {
		log.Printf("Could not persist data: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, fmt.Sprintf("Could not persist data: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Success! Report saved: %+v", report)
	w.WriteHeader(http.StatusOK)
}
