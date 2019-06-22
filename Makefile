.PHONY: all api frontend run clean

all: app frontend

app:
	go build -o bin/epitome github.com/arturovm/epitome/cmd/epitome
	mkdir -p bin/.epitome
	cp -r migrations bin/

run: bin/epitome
	./bin/epitome --debug

clean:
	rm -rf bin
