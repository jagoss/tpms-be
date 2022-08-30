# syntax=docker/dockerfile:1

# Select golang version
FROM golang:1.18 AS build

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir /imgs
ADD . /imgs
# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go .


RUN go build -o /tpms-be ./src/api/main.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /tpms-be /tpms-be
COPY serviceAccountKey.json .

EXPOSE 8080

USER root:root

ENTRYPOINT ["/tpms-be"]
