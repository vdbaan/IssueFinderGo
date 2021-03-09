hello:
	echo "Hello"

compile:
	go build -o bin/main app/server/issuefinder.go -o issuefinder

generate:
	go generate ./...
	cd
	java -jar antlr-4.9-complete.jar -Dlanguage=Go -visitor -o ./infra/filters/parser ./infra/filters/Predicate.g4

build: generate compile

run:
	go run app/server/issuefinder.go

all: hello build
