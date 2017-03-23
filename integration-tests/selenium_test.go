package gamescore_test

import (
	"net"
	"os"
	"os/exec"
	"testing"

	"github.com/tebeka/selenium"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&GamescoreIntegrationTestSuite{})

type GamescoreIntegrationTestSuite struct {
	gamescore *exec.Cmd
	wd        selenium.WebDriver

	sel *selenium.Service
}

func (g *GamescoreIntegrationTestSuite) SetUpSuite(c *C) {
	sel, err := selenium.NewSeleniumService("./selenium-server-standalone-3.0.1.jar", 4444)
	c.Assert(err, IsNil)
	g.sel = sel
}

func (g *GamescoreIntegrationTestSuite) TearDownSuite(c *C) {
	g.sel.Stop()
}

func (g *GamescoreIntegrationTestSuite) SetUpTest(c *C) {
	g.gamescore = exec.Command("./gamescore")
	g.gamescore.Dir = ".."
	g.gamescore.Env = append(os.Environ(), "GAMESCORE_SKIP_LAUNCH_BROWSER=1")
	err := g.gamescore.Start()
	c.Assert(err, IsNil)

	for i := 0; i < 60; i++ {
		c, err := net.Dial("localhost", "8080")
		if err == nil {
			c.Close()
			break
		}
	}

	caps := selenium.Capabilities{
		/*
			"browserName":            "firefox",
			"webdriver.gecko.driver": "./geckodriver",
		*/
		"browserName":       "htmlunit",
		"javascriptEnabled": true,
	}
	wd, err := selenium.NewRemote(caps, "")
	c.Assert(err, IsNil)
	g.wd = wd
}

func (g *GamescoreIntegrationTestSuite) TearDownTest(c *C) {
	g.wd.Quit()
	g.gamescore.Process.Kill()
}

func (g *GamescoreIntegrationTestSuite) TestEditPage(c *C) {
	err := g.wd.Get("http://localhost:8080/score_edit.html")
	c.Assert(err, IsNil)
	g.testCommonElements(c)
}

func (g *GamescoreIntegrationTestSuite) TestDisplayPage(c *C) {
	err := g.wd.Get("http://localhost:8080/score.html")
	c.Assert(err, IsNil)
	g.testCommonElements(c)
}

func (g *GamescoreIntegrationTestSuite) testCommonElements(c *C) {
	for _, elm := range []string{"score_team1", "score_team2", "time"} {
		we, err := g.wd.FindElement(selenium.ByID, elm)
		c.Assert(err, IsNil)
		displayed, err := we.IsDisplayed()
		c.Assert(err, IsNil)
		c.Check(displayed, Equals, true, Commentf("elm %v not displayed", elm))
	}
}
