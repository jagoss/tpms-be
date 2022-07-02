# syntax=docker/dockerfile:1

FROM golang:1.18

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN go get github.com/jagoss/tpms-be/src/api
RUN go /build && git clone https://github.com/jagoss/tpms-be.git

RUN cd /build/tpms/main && go build

EXPOSE 8080

ENTRYPOINT ["/build/tmps/main/mian"]
