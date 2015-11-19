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

type StateChange struct {
	TeamId      string
	ScoreChange int
}

var currentGame Game

func tick() {
	for {
		if !currentGame.Running {
			time.Sleep(100 * time.Millisecond)
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

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return body, nil
}

/*
 create via:
 $ curl -i -H "Content-Type: application/json" -X POST -d '{"Team1":{"Name":"Foo","Goals":2},"Team2":{"Name":"Bar","Goals":7},"TimeLeft":120,"Half":1,"Running":true}' http://localhost:8080/api/1/game
*/
func create(w http.ResponseWriter, r *http.Request) {
	var newGame Game
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, &newGame); err != nil {
		log.Printf("body %q resulted in %s", body, err)
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

/*
 change via:
 $ curl -i -H "Content-Type: application/json" -X PUT -d '{"TeamId": "team1", "ScoreChange": 1}' http://localhost:8080/api/1/game
*/
func stateChange(w http.ResponseWriter, r *http.Request) {
	var newState StateChange
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, &newState); err != nil {
		log.Printf("body %q resulted in %s", body, err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch newState.TeamId {
	case "team1":
		currentGame.Team1.Goals += newState.ScoreChange
	case "team2":
		currentGame.Team2.Goals += newState.ScoreChange
	default:
		log.Printf("got invalid stateChange %s %q", newState, body)
		http.Error(w, "meep", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/1/game", status).Methods("GET")
	r.HandleFunc("/api/1/game", create).Methods("POST")
	r.HandleFunc("/api/1/game", stateChange).Methods("PUT")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// ensure our timer is runnning
	go tick()

	listen := "localhost:8080"
	fmt.Printf("create at http://%s/create.html\n", listen)
	fmt.Printf("status at http://%s/score.html\n", listen)
	fmt.Printf("edit score at http://%s/score_edit.html\n", listen)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(listen, nil))
}
