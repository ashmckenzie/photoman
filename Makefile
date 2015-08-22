.PHONY: all clean
.DEFAULT_GOAL := all

deps:
	go get ./...

update:
	go get -u ./...

all: deps
	go build -o photoman

clean:
	rm -f photoman
