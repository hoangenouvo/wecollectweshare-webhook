FROM golang:1.13-alpine as dev

RUN apk add --no-cache make git curl build-base

COPY . /app/

WORKDIR /app

RUN go build -o ./build/wcws wcws

EXPOSE 1323

ENTRYPOINT ["/app/build/wcws"]

