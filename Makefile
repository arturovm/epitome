.PHONY: all api frontend run clean

all: app frontend

app: vendor
	go build -o bin/epitome github.com/arturovm/epitome/cmd/epitome
	mkdir -p bin/.epitome
	cp -r migrations bin/

vendor: Gopkg.toml Gopkg.lock
	dep ensure

run: bin/epitome
	./bin/epitome --debug

clean:
	rm -rf bin