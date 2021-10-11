package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PullRequest struct {
	Title    string
	Number   int
	Html_url string
	User     User `json:"user"`
	Body     string
}

type User struct {
	Login      string `json:"login"`
	Html_url   string `json:"html_url"`
	Avatar_url string `json:"avatar_url"`
}

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	http.HandleFunc("/", Index)
	http.Handle("/public/", http.StripPrefix("/public", fileServer))
	log.Println("starting webserver on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Fetch open pull requests
	resp, err := http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls")
	if err != nil {
		fmt.Printf("Error fetching github information, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}
	// Read http response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	respString := []byte(string(body))
	// Create empty array of pull request structs
	var prs []PullRequest

	// Parse json automagically into pull request array (pointers are neat)
	err = json.Unmarshal(respString, &prs)
	if err != nil {
		fmt.Printf("Error parsing json, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	w.Write([]byte(fmt.Sprintf("nice, %d prs", len(prs))))
}
