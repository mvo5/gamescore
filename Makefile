all: check build buildwin

deps:
	go get github.com/gorilla/mux gopkg.in/check.v1

check: deps
	go test

build: deps main.go main_test.go 
	go build

buildwin:
	GOOS=windows GOARCH=386 go build -o gamescore.exe

clean:
	rm -f gamescore gamescore.exe

dist: build buildwin
	(cd ..; zip -r gamescore-$(shell date '+%Y%m%d-%H%M').zip gamescore)

.PHONY: deps check build buildwin clean dist
