FROM golang:1.22.2 AS build-stage

WORKDIR /app
COPY . .

RUN go mod download


WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o lcode

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest AS build-release-stage

RUN apk update
RUN apk upgrade
RUN apk add ffmpeg

WORKDIR /

COPY --from=build-stage /app/cmd/lcode /go/bin/goose ./

COPY --from=build-stage /app/internal/infra/database/migrations ./migrations

ENV GOOSE_MIGRATION_DIR=/migrations

ENV GOOSE_DRIVER=postgres

ENV GIN_MODE=release

CMD ["./goose", "$PATH_DB", "up"]

CMD ["./lcode"]