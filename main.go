package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	router := CreateRouter()

	ReadBizTrips(&BTs)

	fmt.Println(BTs)

	log.Fatal(http.ListenAndServe(":8080", router))

}
