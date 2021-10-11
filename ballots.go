package main

import (
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public", fileServer))
	log.Println("starting webserver on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
