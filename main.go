package main

import (
	"log"
	"net/http"
	"to-do-list/router"
)

func main() {
	r := router.Router()
	log.Println("Starting server on the port :8084...")
	//Start server on the port 8084
	log.Fatal(http.ListenAndServe(":8084", r))
}
