package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&GamescoreTestSuite{})

type GamescoreTestSuite struct {
	r       *mux.Router
	respRec *httptest.ResponseRecorder
}

func (g *GamescoreTestSuite) SetUpTest(c *C) {
	g.r = makeRouter()
}

func decodeGameJson(c *C, buf []byte) *Game {
	var game Game
	err := json.Unmarshal(buf, &game)
	c.Assert(err, IsNil)

	return &game
}

func (g *GamescoreTestSuite) getStatus(c *C) *Game {
	req, err := http.NewRequest("GET", "/api/1/status", nil)
	c.Assert(err, IsNil)

	respRec := httptest.NewRecorder()
	g.r.ServeHTTP(respRec, req)
	c.Assert(respRec.Code, Equals, 200)

	return decodeGameJson(c, respRec.Body.Bytes())
}

func (g *GamescoreTestSuite) TestGetStatusTrivial(c *C) {
	currentGame = Game{}
	game := g.getStatus(c)
	c.Assert(game, DeepEquals, &Game{})
}

var testGame = &Game{
	Team1: Team{
		Name:  "Foo",
		Goals: 2,
	},
	Team2: Team{
		Name:  "Bar",
		Goals: 7,
	},
	TimeLeft: 120 * time.Second,
	Countdown: &Countdown{},
	Half:     1,
	Running:  false,
}

func (g *GamescoreTestSuite) TestCreateGame(c *C) {
	buf, err := json.Marshal(testGame)
	c.Assert(err, IsNil)

	req, err := http.NewRequest("POST", "/api/1/create", bytes.NewBuffer(buf))
	c.Assert(err, IsNil)

	respRec := httptest.NewRecorder()
	g.r.ServeHTTP(respRec, req)
	c.Assert(respRec.Code, Equals, 201)

	game := g.getStatus(c)
	c.Assert(game, DeepEquals, testGame)
}

func (g *GamescoreTestSuite) TestChangeState(c *C) {
	currentGame = *testGame

	buf, err := json.Marshal(StateChange{
		TeamId:      "team1",
		ScoreChange: 1,
	})
	c.Assert(err, IsNil)
	req, err := http.NewRequest("POST", "/api/1/changeState", bytes.NewBuffer(buf))
	c.Assert(err, IsNil)

	respRec := httptest.NewRecorder()
	g.r.ServeHTTP(respRec, req)
	c.Assert(respRec.Code, Equals, 201)

	c.Assert(currentGame.Team1.Goals, Equals, testGame.Team1.Goals+1)
}

func (g *GamescoreTestSuite) TestTick(c *C) {
	currentGame = *testGame
	currentGame.Countdown.Set(120*time.Second)
	currentGame.Countdown.Start()

	tickOnce()
	c.Assert(currentGame.Countdown.String(), Equals, "01:59")
	for i := 0; i < 11; i++ {
		tickOnce()
	}
	c.Assert(currentGame.Countdown.String(), Equals, "01:58")
}

