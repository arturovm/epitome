.PHONY: all api frontend run clean

all: app frontend

app: bin/epitome
	mkdir -p bin/.epitome
	cp -r migrations bin/

bin/epitome:
	go build -o bin/epitome github.com/arturovm/epitome/cmd/epitome

run: bin/epitome
	./bin/epitome --debug

clean:
	rm -rf bin
