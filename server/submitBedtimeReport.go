package hafenhause

// import (
// 	"log"
// 	"net/http"
// 	"strings"
// )

// // SubmitBedtimeReport receives a bedtimeReport and saves it to a FireStore
// func SubmitBedtimeReport(w http.ResponseWriter, r *http.Request) {
// 	report, err := bedtime.NewReportFromRequest(r)
// 	if err != nil {
// 		log.Printf("decode error: %v", err)
// 		http.Error(w, "Could not decode request", http.StatusBadRequest)
// 		return
// 	}

// 	if errs := report.Validate(); len(errs) > 0 {
// 		http.Error(w, strings.Join(errs, "; "), http.StatusBadRequest)
// 		return
// 	}

// 	if err := report.Set(); err != nil {
// 		log.Printf("Could not save request %+v with error %+v", r, err)
// 		code := http.StatusInternalServerError
// 		http.Error(w, http.StatusText(code), code)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "text/plain")
// }
