package main

import (
	"goji.io"
	"goji.io/pat"
	"gopkg.in/mgo.v2"
	"net/http"
)

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("whereami").C("trips")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		//log this
		panic(err)
	}
}

func main() {

	msesh, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer msesh.Close()

	// TODO: What is Monotonic, should this rather be "Strong"?
	msesh.SetMode(mgo.Monotonic, true)
	ensureIndex(msesh)

	// clean the DB / trips collection
	//err = msesh.DB("whereami").DropDatabase()
	//msesh.DB("whereami").C("trips").RemoveAll()RemoveAll(bson.M{})

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/"), calendarHandler)
	mux.HandleFunc(pat.Get("/next"), calendarHandlerNextMonth)
	mux.HandleFunc(pat.Get("/prev"), calendarHandlerPrevMonth)

	mux.HandleFunc(pat.Get("/trips"), allTripsHandler(msesh))
	mux.HandleFunc(pat.Post("/trips"), addTripHandler(msesh))
	mux.HandleFunc(pat.Get("/trips/:id"), tripByIDHandler(msesh))
	mux.HandleFunc(pat.Put("/trips/:id"), updateTripByIDHandler(msesh))
	mux.HandleFunc(pat.Delete("/trips/:id"), deleteTripByIDHandler(msesh))

	http.ListenAndServe("localhost:8080", mux)

}

// ToDo: Add a Logging Wrapper to the Handlers
/*
func Logging(in http.Handler, pr *RESTRequestT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		in.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s\t%s", r.Method, r.RequestURI, pr.HttpMethod, pr.Pattern, time.Since(start))
	})
}*/
