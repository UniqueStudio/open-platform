FROM golang:1.14-alpine
LABEL maintainer="fredliang"

RUN apk --no-cache add tzdata git ca-certificates && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64  go build -tags=jsoniter ./main.go
ENV GIN_MODE release

RUN chmod +x ./main
CMD ["./main"]