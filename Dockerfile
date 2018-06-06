FROM alpine
LABEL maintainer="fredliang"

RUN apk --no-cache add tzdata  ca-certificates && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
ENV GIN_MODE release
ADD config.deploy.yml /app/config.deploy.yml
ADD main /app/main
RUN chmod +x ./main
CMD ["./main"]