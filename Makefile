all: check build buildwin buildmac

deps:
	go get github.com/gorilla/mux gopkg.in/check.v1
	make -C static

check: deps
	go build
	go test
	(cd integration-tests; ./setup.sh; go get -t . ; go test)

build: deps main.go main_test.go 
	go build

buildwin:
	GOOS=windows GOARCH=386 go build -o gamescore.exe

buildmac:
	GOOS=darwin GOARCH=amd64 go build -o gamescore.macos

clean:
	rm -f gamescore gamescore.exe gamescore.macos

dist: build buildwin buildmac
	(cd ..; zip --exclude '*parts*' --exclude '*prime*' --exclude '*integration-tests*' --exclude '*.git*' -r gamescore-$(shell date '+%Y%m%d-%H%M').zip gamescore)

.PHONY: deps check build buildwin clean dist
