# syntax=docker/dockerfile:1

# Select golang version
FROM golang:1.18 AS build

RUN mkdir /app
ADD . /app
WORKDIR /app

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

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/tpms-be"]
