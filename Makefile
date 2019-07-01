PORT=8000

all: clean amalgam run

amalgam:
	go build -o amalgam ./cmd/amalgam/main.go

run:
	PORT=$(PORT) ./amalgam

clean:
	rm amalgam || :
