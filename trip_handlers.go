package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"goji.io/pat"
	"time"
)

type BizTripT struct {
	ID    int       `json:"id"`
	Start time.Time `json:"start-of-trip"`
	End   time.Time `json:"end-of-trip"`
	Dest  string    `json:"destination"`
}

var (
	BTs          []BizTripT
	//latestIDused int = 0
)

func (p BizTripT) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func allTripsHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("whereami").C("trips")
		var BTs []BizTripT
		err := c.Find(bson.M{}).All(&BTs)
		if err != nil {
			formatError(w, http.StatusInternalServerError, "Database error, Find().All()")
			//log this
			panic(err) //something wrong with the DB, could keep going but maybe something is wrong
		}

		respBody, err := json.MarshalIndent(BTs, "", "  ")
		if err != nil {
			formatError(w, http.StatusInternalServerError, "Marshalling response failed")
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}

func tripByIDHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		idint, err := strconv.Atoi(pat.Param(r, "isbn"))
		if err != nil {
			// business trip ID is not a number
			formatError(w, http.StatusUnprocessableEntity, "Requested ID not a Number")
			return //keep going, problem in the request formatting only
		}

		c := session.DB("whereami").C("trips")

		var trip BizTripT

		err = c.Find(bson.M{"id": idint}).One(&trip)
		if err != nil {
			// log this "Database error, Find().One()"
			panic(err) //something is wrong
		}

		if trip.ID == 0 {
			formatError(w, http.StatusNotFound, "Trip not found")
			return //all is well, keep going
		}

		resp, err := json.MarshalIndent(trip, "", "  ")
		if err != nil {
			// write something before panic
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func formatError(w http.ResponseWriter, errcode int, errmsg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errcode)
	fmt.Fprintf(w, errmsg)
}

func addTripHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var trip BizTripT
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&trip)
		if err != nil {
			formatError(w, http.StatusBadRequest,"Parsing body failed")
			return //keep running no need to panic and hang oneself up, http request was answered
		}

		c := session.DB("whereami").C("trips")

		err = c.Insert(&trip)
		if err != nil {
			//duplicate
			if mgo.IsDup(err) {
				formatError(w, http.StatusBadRequest, "Trip with this ID already exists",)
				return //keep going
			}

			//todo, here we should log something
			panic(err) //there is some problem with the database
		}

		//happy path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", fmt.Sprintf("%s/%d", r.URL.Path, trip.ID))
	}
}

func updateTripByIDHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		idint, err := strconv.Atoi(pat.Param(r, "isbn"))
		if err != nil {
			// business trip ID is not a number
			formatError(w, http.StatusUnprocessableEntity, "Requested ID not a Number")
			return //keep going, problem in the request formatting only
		}

		var trip BizTripT
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&trip)
		if err != nil {
			formatError(w, http.StatusBadRequest, "cant decode JSON from request body")
			return
		}

		c := session.DB("whereami").C("trips")

		err = c.Update(bson.M{"id": idint}, &trip)
		if err != nil {
			if err == mgo.ErrNotFound {
				formatError(w, http.StatusNotFound, "Trip with this ID not found" )
				return
			}
			formatError(w, http.StatusInternalServerError, "Database error, Update()")
			// write some record of this
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteTripByIDHandler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		idint, err := strconv.Atoi(pat.Param(r, "isbn"))
		if err != nil {
			// business trip ID is not a number
			formatError(w, http.StatusUnprocessableEntity, "Requested ID not a Number")
			return //keep going, problem in the request formatting only
		}

		c := session.DB("whereami").C("trips")

		err = c.Remove(bson.M{"id": idint})
		if err != nil {
			if err == mgo.ErrNotFound {
				formatError(w, http.StatusNotFound, "Trip with this ID not found")
				return
			}

			formatError(w, http.StatusInternalServerError, "Database error, Remove()")
			//log this problem
			return
		}

		w.WriteHeader(http.StatusNoContent)
		//todo: should there be a body in this case of response?
	}
}
