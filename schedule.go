package main

// TODO: einradhockeyliga specific, think about generalizing
type Schedule struct {
	Games []struct {
		ID           string `xml:"spiel_id"`
		Team1        string `xml:"teamname_a"`
		Team2        string `xml:"teamname_b"`
		LenHalftimes int    `xml:"halbzeit_laenge"`
		NrHalftimes  int    `xml:"anzahl_halbzeiten"`
	} `xml:"spiel"`
}
