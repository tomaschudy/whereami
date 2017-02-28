package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type DayInCalT struct {
	Number int
	OnTrip bool
}

type Page struct {
	Month string
	Year int
	Day []DayInCalT
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}

	//fmt.Println(currentMonth.Month().String())

	monthstart := time.Date(currentMonth.Year(), currentMonth.Month(), 1 ,
				0, 0, 0,
				0, currentMonth.Location())
	monthend := time.Date(currentMonth.Year(), currentMonth.Month(), 1 ,
				23, 59, 59,
				0, currentMonth.Location())
	monthend = monthend.AddDate (0, 1, -1)

	//fmt.Println(monthstart, monthend)

	firstday := monthstart
	for (firstday.Weekday() > 0) {
		firstday = firstday.AddDate(0,0, -1)
	}

	lastday := monthend
	for (lastday.Weekday() < 6) {
		lastday = lastday.AddDate(0,0,1)
	}

	var pd []DayInCalT

	d := firstday //move time to 23:58 and it will be OK

	for lastday.After(d) {
		pd = append(pd, DayInCalT{d.Day(), IsOnBusTrip(d)})
		d = d.AddDate(0, 0, 1)
	}

	p := &Page{ currentMonth.Month().String(), currentMonth.Year(), pd }

	renderTemplate(w, "Calendar", p)
}

func nextMonthHandler(w http.ResponseWriter, r *http.Request) {

	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}

	currentMonth = currentMonth.AddDate(0,1, 0)

	http.Redirect(w, r, "/", http.StatusFound)
}

func prevMonthHandler(w http.ResponseWriter, r *http.Request) {

	if currentMonth.IsZero() {
		currentMonth = time.Now()
	}

	currentMonth = currentMonth.AddDate(0, -1, 0)

	http.Redirect(w, r, "/", http.StatusFound)
}

func IsOnBusTrip(d time.Time) bool {
	for t := range BTs {
		if (d.After(BTs[t].Start) && d.Before(BTs[t].End)) { return true }
	}
	return false
}

func editBusinessTripsHandler(w http.ResponseWriter, r *http.Request) {
	var trips []BizTripT
	readBizTripz(&trips)
	fmt.Println(len(trips), cap(trips), trips)

	t, err := template.ParseFiles(  "EditBusinessTrips.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	return
	}

	err = t.Execute(w, &trips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func enterBusTripHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles(  "EnterNewBusinessTrip.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, []int{1,2,3})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	r.ParseForm()
	fmt.Println("starttime:", r.Form["starttime"])
	fmt.Println("endtime:", r.Form["endtime"])
	fmt.Println("destination:", r.Form["destination"])

	r.ParseForm()
	fmt.Println("starttime:", r.Form["starttime"])
	fmt.Println("endtime:", r.Form["endtime"])
	fmt.Println("destination:", r.Form["destination"])
	/*
		layout := "2006-01-02T15:04"
		s, err := time.Parse( layout, r.Form["starttime"][0])
		if err != nil {
			fmt.Println(err)
		}

		e, err := time.Parse( layout, r.Form["endtime"][0])
		if err != nil {
			fmt.Println(err)
		}
	/*
		fmt.Println(cap(BTs), len(BTs))

		nt := BizTripT{1, s, e, r.Form["destination"][0] }
		BTs[1] = nt


		//nt.saveBizTripz([]BizTripT(nt))
		fmt.Println(toJson(BTs[1]), len(BTs))
		*/


	http.Redirect(w, r, "/", http.StatusFound)
}

type BizTripT struct {
	ID uint `json:"id"`
	Start time.Time `json:"start-of-trip"`
	End time.Time `json:"end-of-trip"`
	Dest string `json:"destination"`
}

func (p BizTripT) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

func readBizTripz(p *[]BizTripT) {
	filedata, err := ioutil.ReadFile("./biztripz.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(filedata, p)

	//fmt.Println(c, err)
	return
}

func saveBizTripz(s *[]BizTripT) {
	for i:=0; i<len(*s); i++ {
		fmt.Println(*s)
		err := ioutil.WriteFile("biztripz.json", []byte((*s)[i].toString()), 0600)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

var BTs []BizTripT
var currentMonth time.Time

func main() {
	readBizTripz(&BTs)
	//fmt.Println(len(BTs), cap(BTs), BTs)
	//saveBizTripz(&BTs)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/NextMonth/", nextMonthHandler)
	http.HandleFunc("/PrevMonth/", prevMonthHandler)
	http.HandleFunc("/EditBusinessTrips/", editBusinessTripsHandler)
	http.ListenAndServe(":8080", nil)
}