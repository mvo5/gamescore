package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

type Team struct {
	Name  string
	Goals int
}

type Game struct {
	Team1 Team
	Team2 Team

	// XXX: cleanup and handle all of this via "Countdown"
	TimeLeft time.Duration
	TimeStr  string

	Half    int
	Running bool

	Countdown *Countdown
}

type StateChange struct {
	TeamId        string `json:",omitempty"`
	ScoreChange   int    `json:",omitempty"`
	ToggleTimeout bool   `json:",omitempty"`
}

// XXX: mutex!
var (
	currentGame     Game
	currentSchedule Schedule
)

func init() {
	// init with some example data
	currentGame = Game{
		Team1: Team{
			Name: "Home",
		},
		Team2: Team{
			Name: "Guest",
		},
		Countdown: &Countdown{},
	}
}

func tickOnce() {
	if currentGame.Countdown.TimeLeft() <= 0 {
		currentGame.Running = false
		currentGame.TimeStr = "00:00"
		currentGame.TimeLeft = 0
	} else {
		currentGame.TimeStr = currentGame.Countdown.String()
		currentGame.TimeLeft = currentGame.Countdown.TimeLeft()
	}
	time.Sleep(100 * time.Millisecond)
}

func tick() {
	for {
		if currentGame.Running {
			currentGame.Countdown.Start()
		} else {
			currentGame.Countdown.Stop()
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

	newGame.Countdown = NewCountdown(newGame.TimeLeft)
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
	}

	if newState.ToggleTimeout {
		currentGame.Running = !currentGame.Running
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

/*
 create via:
 $ curl -i -H "Content-Type: application/xml" -X POST -d '@test-data/einradhockey1.xml' http://localhost:8080/api/1/schedule
*/
func postSchedule(w http.ResponseWriter, r *http.Request) {
	var newSchedule Schedule
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	if err := xml.Unmarshal(body, &newSchedule); err != nil {
		log.Printf("body %q resulted in %s", body, err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentSchedule = newSchedule
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

/* show via:
 $ curl http://localhost:8080/api/1/schedule
 curl http://localhost:8080/api/1/schedule
{"Games":[{"ID":"1","Team1":"Crazy Ducks","Team2":"Quer-durchs-Land*","LenHalftimes":15,"NrHalftimes":2},{"ID":"2","Team1":"Black Stars","Team2":"MJC Trier - Die Römer","LenHalftimes":15,"NrHalftimes":2}]}
*/
func getSchedule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	outJSON, err := json.Marshal(currentSchedule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", outJSON)
}

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/1/status", status).Methods("GET")
	r.HandleFunc("/api/1/create", create).Methods("POST")
	r.HandleFunc("/api/1/changeState", stateChange).Methods("POST")

	r.HandleFunc("/api/1/schedule", postSchedule).Methods("POST")
	r.HandleFunc("/api/1/schedule", getSchedule).Methods("GET")

	prefix := os.Getenv("SNAP")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(filepath.Join(prefix, "./static/"))))

	return r
}

func launchBrowser(editUrl, statusUrl string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", statusUrl).Start()
		exec.Command("xdg-open", editUrl).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", statusUrl).Start()
		exec.Command("cmd", "/c", "start", editUrl).Start()
	case "darwin":
		exec.Command("open", statusUrl).Start()
		exec.Command("open", editUrl).Start()
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
