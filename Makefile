all: check build buildwin

check:
	go test

build: main.go main_test.go
	go build

buildwin:
	GOOS=windows GOARCH=386 go build -o gamescore.exe

clean:
	rm -f gamescore gamescore.exe
