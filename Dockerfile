# syntax=docker/dockerfile:1

# Select golang version
FROM golang:1.18

RUN mkdir /app
ADD . /app
WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go .

RUN go build -o /src/api/main .

EXPOSE 8080

CMD [ "/tmps-be" ]
