hello:
	@echo "Hello"

compile: generate
	GOOS=windows GOARCH=amd64 go get -v
	GOOS=windows GOARCH=amd64 go build -o bin/issuefinder.exe -ldflags app/server/issuefinder.go

generate:
	go generate ./...
	wget http://www.antlr.org/download/antlr-4.9-complete.jar
	java -jar antlr-4.9-complete.jar -Dlanguage=Go -visitor -o ./infra/filters/parser ./infra/filters/Predicate.g4

build: generate compile

run:
	go run app/server/issuefinder.go

all: hello build

test:
	env GOOS=windows GOARCH=amd64 go test -v ./...

