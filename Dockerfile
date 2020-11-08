FROM golang:1.15.4-alpine as build

COPY . /app
WORKDIR /app

RUN apk update
RUN apk upgrade
RUN apk add --update gcc musl-dev g++

RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o main main.go



FROM alpine:latest
COPY --from=build /app /app
WORKDIR /app

CMD /app/main