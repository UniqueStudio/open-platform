export PROJECT_PATH = $GOPATH/src/family-tree

all: test

env:
	export GIN_MODE=test && \
	swag init

run: env
	go run ./main.go

clean:
	rm -ef *.json docs