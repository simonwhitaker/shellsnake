# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /shellsnake

FROM alpine:3.16

WORKDIR /

COPY --from=build /shellsnake /shellsnake

CMD [ "/shellsnake" ]
