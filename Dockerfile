FROM golang:alpine as builder
ARG PROJECT_ROOT
ARG GIT_TAG_NAME
WORKDIR $PROJECT_ROOT

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -ldflags "-X main.Version=$GIT_TAG_NAME" -o bin/main ./cmd/app

#Meant for building the deployment container
FROM golang:alpine
ARG PROJECT_ROOT

WORKDIR /go

RUN apk update && \
    apk add ca-certificates && \
    apk add --no-cache bash && \
    rm -rf /var/cache/apk/*

COPY --from=builder $PROJECT_ROOT/bin ./
ENTRYPOINT ["./main"]