# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /emojisnake

FROM alpine:3.16

WORKDIR /

COPY --from=build /emojisnake /emojisnake

CMD [ "/emojisnake" ]
