FROM golang:1.22.2 AS build-stage

WORKDIR /app
COPY . .

RUN go mod download


WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o lcode

FROM alpine:latest AS build-release-stage

RUN apk update
RUN apk upgrade
RUN apk add ffmpeg

WORKDIR /

RUN wget https://github.com/pressly/goose/releases/download/v3.19.2/goose_linux_x86_64

COPY --from=build-stage /app/cmd/lcode /app/scripts ./

COPY --from=build-stage /app/docs ./docs

COPY --from=build-stage /app/internal/infra/database/migrations ./migrations


ENV GOOSE_MIGRATION_DIR=/migrations

ENV GOOSE_DRIVER=postgres

ENV GIN_MODE=release

RUN chmod +x ./run.sh

ENTRYPOINT ["/bin/sh", "./run.sh"]