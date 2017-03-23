package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
	TimeStr  string
	Half     int
	Running  bool
}

type StateChange struct {
	TeamId        string `json:",omitempty"`
	ScoreChange   int    `json:",omitempty"`
	ToggleTimeout bool   `json:",omitempty"`
}

var currentGame Game

func init() {
	// init with some example data
	currentGame = Game{
		Team1: Team{
			Name: "Home",
		},
		Team2: Team{
			Name: "Guest",
		},
	}
	// and ensure time is availale
	formatTime()
}

func tickOnce() {
	currentGame.TimeLeft -= time.Duration(100 * time.Millisecond)
	if currentGame.TimeLeft <= 0 {
		currentGame.Running = false
	}

	formatTime()
	time.Sleep(100 * time.Millisecond)
}

func formatTime() {
	// format the time nicely for the javascript
	min := currentGame.TimeLeft / (60 * time.Second)
	sec := currentGame.TimeLeft % (60 * time.Second)
	// sub 1 min gets 100 millisecond resolution
	if min >= 1 {
		sec /= time.Second
		currentGame.TimeStr = fmt.Sprintf("%02d:%02d", min, sec)
	} else {
		sec /= 100 * time.Millisecond
		currentGame.TimeStr = fmt.Sprintf("%02d:%03d", min, sec)
	}
}

func tick() {
	for {
		if !currentGame.Running {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		tickOnce()
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

	newGame.TimeLeft = time.Duration(newGame.TimeLeft)
	currentGame = newGame
	formatTime()

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
	}

	if newState.ToggleTimeout {
		currentGame.Running = !currentGame.Running
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/1/status", status).Methods("GET")
	r.HandleFunc("/api/1/create", create).Methods("POST")
	r.HandleFunc("/api/1/changeState", stateChange).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return r
}

func launchBrowser(editUrl, statusUrl string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", statusUrl).Start()
		exec.Command("xdg-open", editUrl).Start()
	case "windows", "darwin":
		exec.Command("cmd", "/c", "start", statusUrl).Start()
		exec.Command("cmd", "/c", "start", editUrl).Start()
	default:
		fmt.Println("unsupported platform")
	}
}

func main() {
	r := makeRouter()

	// ensure our timer is runnning
	go tick()

	listen := "localhost:8080"
	editUrl := fmt.Sprintf("http://%s/score_edit.html", listen)
	statusUrl := fmt.Sprintf("http://%s/score.html", listen)
	fmt.Printf("edit game at %s\n", editUrl)
	fmt.Printf("display status at %s\n", statusUrl)

	http.Handle("/", r)
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("GAMESCORE_SKIP_LAUNCH_BROWSER") == "" {
		launchBrowser(editUrl, statusUrl)
	}
	log.Fatal(http.Serve(listener, nil))
}
