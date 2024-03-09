package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEinradhockeyligaParseXml(t *testing.T) {
	var schedule Schedule

	_, filename, _, _ := runtime.Caller(0)
	schedulexml, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "test-data/einradhockey1.xml"))
	assert.NoError(t, err)

	err = xml.Unmarshal([]byte(schedulexml), &schedule)
	assert.NoError(t, err)
	assert.Equal(t, "1", schedule.Games[0].ID)
	assert.Equal(t, 15, schedule.Games[0].LenHalftimes)
	assert.Equal(t, 2, schedule.Games[0].NrHalftimes)
	assert.Equal(t, "Crazy Ducks", schedule.Games[0].Team1)
	assert.Equal(t, "Quer-durchs-Land*", schedule.Games[0].Team2)

	assert.Equal(t, 10, len(schedule.Games))
}
