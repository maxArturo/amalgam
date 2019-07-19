PORT=8000

all: clean test amalgam run

amalgam:
	go build -o amalgam ./cmd/amalgam/main.go

run:
	PORT=$(PORT) ./amalgam

clean:
	rm amalgam || :

test: 
	go test ./... --timeout 1s
