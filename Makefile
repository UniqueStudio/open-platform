export PROJECT_PATH = $GOPATH/src/family-tree

all: test

env:
	export GIN_MODE=test

run: env
	go run ./main.go

deploy:
	GOOS=linux GOARCH=amd64  go build -tags=jsoniter ./main.go
	docker build -t registry.cn-hangzhou.aliyuncs.com/fredliang/open-platform  .
	docker push  registry.cn-hangzhou.aliyuncs.com/fredliang/open-platform

clean:
	rm -rf *.json docs