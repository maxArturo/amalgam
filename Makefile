PORT=8000
GO111MODULE=on

all: clean test amalgam run

amalgam: clean
	go build -o amalgam ./cmd/amalgam/main.go

run: amalgam
	PORT=$(PORT) ./amalgam

clean:
	rm amalgam || :

test: 
	go test ./... --timeout 1s
