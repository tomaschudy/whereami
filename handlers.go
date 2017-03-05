package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//error
type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

var currentMonth time.Time

// TODO: This is to return a calendar
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}

	month, err := FormatCalMonth(currentMonth)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(month); err != nil {
		panic(err)
	}

}

func HomeHandlerNext(w http.ResponseWriter, r *http.Request) {
	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}
	currentMonth = currentMonth.AddDate(0, 1, 0)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HomeHandlerPrev(w http.ResponseWriter, r *http.Request) {
	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}
	currentMonth = currentMonth.AddDate(0, -1, 0)
	http.Redirect(w, r, "/", http.StatusFound)
}

func FormatResponse(httpstatus int, errText string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpstatus)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: httpstatus, Text: errText}); err != nil {
		panic(err)
	}
}

func ListAllBusTripsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(BTs); err != nil {
		panic(err)
	}
}

func FindBusTripHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r) //returns a map[string]string
	var inBTID int

	inBTID, err := strconv.Atoi(params["businesstripID"])
	if err != nil {
		// business trip ID is not a number
		FormatResponse(http.StatusUnprocessableEntity, "Requested ID not a Number", w)
		return
	}

	trip, err := FindBizTrip(inBTID)
	if err != nil {
		FormatResponse(http.StatusNotFound, "Not Found", w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(trip); err != nil {
		panic(err)
	}
	return
}

func NewBusTripHandler(w http.ResponseWriter, r *http.Request) {
	var trip BizTripT

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &trip); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := AppendBizTrip(trip)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}

	SaveBizTrips(&BTs)
}
