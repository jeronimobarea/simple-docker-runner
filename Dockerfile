FROM golang:1.20-alpine AS build

WORKDIR /go/src/app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go install -v ./cmd/...

FROM alpine:latest as docker

WORKDIR /app
COPY allowed_images.json ./
COPY --from=build /go/bin/docker docker
CMD ["./docker"]
