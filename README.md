[![Build Status][travis-image]][travis-url]

# gamescore

gamescore is a system to display and edit the gamescore of a sports game. It conatins a small json API to
create the game and update the score plus a html/js frontend to display the stats like time left and goals.

Examle:
```
$ (cd static ; ./get-jquery.sh)
$ go build && ./gamescore
edit game at http://localhost:8080/score_edit.html
display status at http://localhost:8080/score.html
```

## Example use

The way it is desinged is that you have the browser windows with
`score_edit.html` on your laptop display. It is used to control
the time and goals in the match. Use an external monitor to display
the current time and score to the players on the pitch.

## TODO

Add option to select if `score_html.html` is displayed mirrored
or not. If you use an external monitor mirroring is what we want
so that you and the external monitor have the same team names on
the same sides. However if e.g. a beamer is used, things are
different and the teams may not need to be mirrored.

# Integration with einradhockeyliga XML

Get the xml from 
https://einrad.hockey/xml/spielplan?turnier_id=1021
and upload in the button.

Get the "turnier_id" by clicking at "Turnier details" and check the
URL (e.g. https://einrad.hockey/liga/turnier_details?turnier_id=1198)


[travis-image]: https://travis-ci.org/mvo5/gamescore.svg?branch=master
[travis-url]: https://travis-ci.org/mvo5/gamescore
