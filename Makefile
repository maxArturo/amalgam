PORT=8000

all: clean amalgam

amalgam:
	go build ./internal/app/amalgam

run: 
	PORT=$(PORT) ./amalgam

clean:
	rm amalgam
