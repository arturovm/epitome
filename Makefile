.PHONY: all api frontend clean

all: app frontend

app: bin/epitome

bin/epitome: vendor
	go build -o bin/epitome github.com/Arturovm/epitome/cmd/epitome

vendor: Gopkg.toml Gopkg.lock
	dep ensure

clean:
	rm -rf bin