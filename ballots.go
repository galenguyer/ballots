package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	http.HandleFunc("/ballot", Ballot)
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

	// TODO: Should find a way to handle not having any PRS... Otherwise template rendering will just crash
	data := struct {
		Prs []PullRequest
	}{
		Prs: prs,
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Error executing template, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}
}

func Ballot(w http.ResponseWriter, r *http.Request) {
	p, present := r.URL.Query()["pr"]
	if !present {
		w.Write([]byte("Specify a PR"))
		return
	}
	gh_pr := p[0]

	// Fetch pull request diff
	resp, err := http.Get(fmt.Sprintf("https://patch-diff.githubusercontent.com/raw/ComputerScienceHouse/Constitution/pull/%s.diff", gh_pr))
	if err != nil {
		fmt.Printf("Error fetching PR diff, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	diffString := string(body)
	// Fix new lines
	strings.Replace(diffString, `\n`, "\n", -1)

	// Fetch PR data to determine the title
	resp, err = http.Get(fmt.Sprintf("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls/%s.diff", gh_pr))
	if err != nil {
		fmt.Printf("Error fetching PR title, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}
	defer resp.Body.Close()
	titleBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	// Cast response body to usable byte array
	titleString := []byte(string(titleBody))
	var pr PullRequest
	// Parse title with json
	err = json.Unmarshal(titleString, &pr)
	if err != nil {
		fmt.Printf("Error parsing json, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

	// TODO
	// Variable number of pokemon/ballots
	pokemons := getPokemon(100)
	data := struct {
		Pokemons map[int]string
		DiffStr  string
		Pr       PullRequest
	}{
		Pokemons: pokemons,
		DiffStr:  diffString,
		Pr:       pr,
	}
	tmpl := template.Must(template.ParseFiles("templates/ballot.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("Error executing template, %s\n", err)
		w.Write([]byte("Oops, a fucky wucky occured"))
		return
	}

}

func getPokemon(numBallots int) map[int]string {
	// Open pokemon csv file
	pokefile, err := os.Open("./pokemon.csv")
	if err != nil {
		fmt.Printf("Error opening pokemon.csv, %s\n", err)
		return nil
	}

	r := csv.NewReader(bufio.NewReader(pokefile))
	// Off by one as usual
	numBallots = numBallots + 1
	// Create empty pokemon array (this is some weird go call that just works)
	pokemons := make(map[int]string, numBallots)
	for i := 1; i < numBallots; i++ {
		pokemon, err := r.Read()
		if err == io.EOF {
			break
		}
		pokemons[i] = pokemon[1]
	}
	// Closing file is important
	pokefile.Close()
	return pokemons
}
