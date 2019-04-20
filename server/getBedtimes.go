package hafenhause

// import (
// 	"encoding/json"
// 	"net/http"
// )

// // GetBedtimes gets Bran and Malcolm's current bedtimes
// func GetBedtimes(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET")
// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

<<<<<<< Updated upstream
// GetBedtimes gets Bran and Malcolm's current bedtimes
func GetBedtimes(w http.ResponseWriter, r *http.Request) {
	bedtimes, err := bedtime.GetAllBedtimes()
=======
// 	bedtimes, err := bedtime.GetAll()

// 	bedtimesJSON, err := json.Marshal(bedtimes)
// 	if err != nil {
// 		panic(err)
// 	}
>>>>>>> Stashed changes

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(bedtimesJSON)
// }
