export PROJECT_PATH = $GOPATH/src/family-tree

all: test

env:
	export GIN_MODE=test

run: env
	go run ./main.go

clean:
	rm -rf *.json docs