.PHONY: all clean

deps:
	go get -u ./...

all: deps
	go build -o photoman

clean:
	rm -f photoman
