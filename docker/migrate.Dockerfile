FROM golang:1.25-alpine

RUN apk add --no-cache ca-certificates bash

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY migrations ./migrations

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=/app/migrations

ENTRYPOINT ["goose"]
CMD ["up"]

