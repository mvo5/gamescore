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

[travis-image]: https://travis-ci.org/mvo5/gamescore.svg?branch=master
[travis-url]: https://travis-ci.org/mvo5/gamescore
