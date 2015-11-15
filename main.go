package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Team struct {
	Name  string
	Goals int
}

type Game struct {
	Team1    Team
	Team2    Team
	TimeLeft time.Duration
	Half     int
	Running  bool
}

var currentGame Game

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/1/game", status).Methods("GET")
	router.HandleFunc("/api/1/game", create).Methods("POST")

	// ensure our timer is runnning
	go tick()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func tick() {
	for {
		if !currentGame.Running {
			continue
		}
		currentGame.TimeLeft -= time.Duration(100 * time.Millisecond)
		if currentGame.TimeLeft <= 0 {
			currentGame.Running = false
		}
		time.Sleep(100 * time.Millisecond)
	}
}

/* show via:
 $ curl http://localhost:8080/api/1/game
{"Team1":{"Name":"Foo","Goals":2},"Team2":{"Name":"Bar","Goals":7},"TimeLeft":86500000000,"Half":1,"Running":true}
*/
func status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	outJSON, err := json.Marshal(currentGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", outJSON)
}

/*
 create via:
 $ curl -i -H "Content-Type: application/json" -X POST -d '{"Team1":{"Name":"Foo","Goals":2},"Team2":{"Name":"Bar","Goals":7},"TimeLeft":120,"Half":1,"Running":true}' http://localhost:8080/api/1/game
*/
func create(w http.ResponseWriter, r *http.Request) {
	var newGame Game
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &newGame); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGame.TimeLeft = time.Duration(newGame.TimeLeft * time.Second)

	currentGame = newGame
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}