FROM docker.io/library/golang:1.17-alpine AS builder

WORKDIR /ballots/
COPY . .
RUN go build -v

FROM alpine:latest
COPY --from=builder /ballots/ballots /ballots
COPY ./templates /templates
COPY ./public /public
COPY ./pokemon.csv /pokemon.csv
ENTRYPOINT /ballots