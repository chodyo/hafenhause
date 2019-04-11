package hafenhause

import (
	"encoding/json"
	"net/http"

	"github.com/chodyo/hafenhause/bedtime"
)

// GetBedtimes gets Bran and Malcolm's current bedtimes
func GetBedtimes(w http.ResponseWriter, r *http.Request) {
	bedtimes, err := bedtime.GetAllBedtimes()

	bedtimesJSON, err := json.Marshal(bedtimes)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bedtimesJSON)
}
