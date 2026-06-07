GH_NAME:=create-iteration

build:
	go build -o gh-${GH_NAME} main.go

install: build
	gh extension remove ${GH_NAME} || echo
	gh extension install .

start: install
	gh ${GH_NAME} -project-url=https://github.com/users/notomo/projects/1 -field=Iteration -count=3 -duration=7 -dry-run -log=/dev/stdout

test:
	go test -v ./...

help:
	go run main.go -h
