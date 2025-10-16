package main

import (
	"fmt"
	"log"
	"net/http"

	gosltimetable "github.com/alexdriaguine/go-sl-time-table/internal"
	"github.com/alexdriaguine/go-sl-time-table/internal/sl_api"
)

func main() {
	fmt.Println("SL Time Table server!")

	// move to a server
	port := ":3000"

	slClient := sl_api.NewDefaultSLApi()
	router, err := gosltimetable.NewRouter(slClient)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("started server on port %s\n", port)
	err = http.ListenAndServe(port, router)

	log.Fatal(err)

}
