FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY="https://goproxy.cn,direct" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -tags=jsoniter -ldflags="-w -s" main.go

FROM xylonx/cn-ubuntu:latest AS prod
ARG PROJECT_NAME=open-platform

LABEL owner="UniqueStudio"

WORKDIR /opt/${PROJECT_NAME}

COPY --from=builder /app/main ./${PROJECT_NAME}

RUN echo "./${PROJECT_NAME} -c settings.yaml" > run.sh &&\
    chmod -R 755 /opt/${PROJECT_NAME}

EXPOSE 5000

CMD ./run.sh
#FROM alpine
#LABEL maintainer="fredliang"
#
#RUN apk --no-cache add tzdata  ca-certificates && \
#    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#    echo "Asia/Shanghai" > /etc/timezone
#
#WORKDIR /app
#ENV GIN_MODE release
#ADD static /app/static
#ADD docs /app/docs
#
#ADD main /app/main
#RUN chmod +x ./main
#CMD ["./main"]