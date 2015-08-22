.PHONY: all clean

all:
	go build -o photoman

clean:
	rm -f photoman
