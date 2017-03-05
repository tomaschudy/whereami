package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	latestIDused int = 0
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

func ReadBizTrips(p *[]BizTripT) {
	filedata, err := ioutil.ReadFile("./biztripz.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(filedata, p)
	latestIDused = len(*p)

	fmt.Println(err)
	return
}

//TODO: Clean up this "fancy marshalling"
func SaveBizTrips(s *[]BizTripT) {
	bytes, err := json.Marshal(s)
	err = ioutil.WriteFile("biztripz.json", bytes, 0600)
	if err != nil {
		panic(err)
	}
}

func FindBizTrip(id int) (BizTripT, error) {
	for _, t := range BTs {
		if t.ID == id {
			return t, nil
		}
	}

	// return empty BizTripT if not found
	return BizTripT{}, errors.New("not found")
}

//TODO: race condition safety?
func AppendBizTrip(t BizTripT) BizTripT {
	latestIDused++
	t.ID = latestIDused
	BTs = append(BTs, t)
	return t
}

func DeleteBizTrip(id int) error {
	for i, t := range BTs {
		if t.ID == id {
			BTs = append(BTs[:i], BTs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("DeleteBizTrip; Business Trip ID: %d not found.", id)
}
