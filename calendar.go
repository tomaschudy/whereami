package main

import (
	"fmt"
	"time"
)

type DayT struct {
	Number int  `json:"date"`
	OnTrip bool `json:"on-trip"`
}

type MonthT struct {
	Month string `json:"month"`
	Year  int    `json:"year"`
	Days  []DayT `json:"days"`
}

func FormatCalMonth(t time.Time) (*MonthT, error) {

	monthstart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	monthend := time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 0, t.Location())
	monthend = monthend.AddDate(0, 1, -1)

	fmt.Println(monthstart, monthend)

	firstday := monthstart
	for firstday.Weekday() > 0 {
		firstday = firstday.AddDate(0, 0, -1)
	}

	lastday := monthend
	for lastday.Weekday() < 6 {
		lastday = lastday.AddDate(0, 0, 1)
	}

	month := MonthT{t.Month().String(), t.Year(), []DayT{}}

	//move time to 23:58 and it will be OK
	d := time.Date(firstday.Year(), firstday.Month(), firstday.Day(), 23, 59, 58, firstday.Nanosecond(), firstday.Location())

	for lastday.After(d) {
		month.Days = append(month.Days, DayT{d.Day(), IsOnBusTrip(d)})
		d = d.AddDate(0, 0, 1)
	}

	return &month, nil
}

//TODO: there is a bug here - false negative for when d is late (after .End on a given day
func IsOnBusTrip(d time.Time) bool {
	for t := range BTs {
		if d.After(BTs[t].Start) && d.Before(BTs[t].End) {
			return true
		}
	}
	return false
}
