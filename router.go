package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

//defines REST request patters and handlers
type RESTRequestT struct {
	HttpMethod  string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type RESTRequestsSliceT []RESTRequestT

var requests = RESTRequestsSliceT{
	// home shows a calendar
	RESTRequestT{"GET", "/", HomeHandler},
	//next month
	RESTRequestT{"GET", "/next", HomeHandlerNext},
	//previous month
	RESTRequestT{"GET", "/prev", HomeHandlerPrev},
	//business trip handling
	RESTRequestT{"GET", "/businesstrips", ListAllBusTripsHandler},
	RESTRequestT{"POST", "/businesstrips", NewBusTripHandler},
	RESTRequestT{"GET", "/businesstrips/{businesstripID}", FindBusTripHandler},
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, req := range requests {
		var handler http.Handler

		handler = req.HandlerFunc
		handler = Logging(handler, &req)

		router.Methods(req.HttpMethod).Path(req.Pattern).Handler(handler)
	}

	return router
}

//logger wrapper
func Logging(in http.Handler, pr *RESTRequestT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		in.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s\t%s", r.Method, r.RequestURI, pr.HttpMethod, pr.Pattern, time.Since(start))
	})
}
