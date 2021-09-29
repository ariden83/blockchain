#Meant for building the deployment container
FROM alpine:3.10.1
ARG PROJECT_ROOT

WORKDIR /go

RUN apk update && \
    apk add ca-certificates && \
    apk add --no-cache bash && \
    rm -rf /var/cache/apk/*

COPY . .
ENTRYPOINT ["./main"]
