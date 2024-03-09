package main

// TODO: einradhockeyliga specific, think about generalizing
type Schedule struct {
	Games []struct {
		ID           string `xml:"spiel_id" json:"id"`
		Team1        string `xml:"teamname_a" json:"team1"`
		Team2        string `xml:"teamname_b" json:"team2"`
		LenHalftimes int    `xml:"halbzeit_laenge" json:"len_half"`
		NrHalftimes  int    `xml:"anzahl_halbzeiten" json:"nr_half"`
	} `xml:"spiel" json:"games"`
}
