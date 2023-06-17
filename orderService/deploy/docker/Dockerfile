FROM golang:1.20 as build

WORKDIR /build

COPY orderService/Makefile .
COPY ./orderService/internal/ /build/internal/
COPY go.mod .
COPY ./orderService/docs/ /build/docs/
COPY ./orderService/cmd/ /build/cmd/

RUN make download prebuild build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/authsvc .

COPY orderService/config/config.yaml orderService/config.yaml

CMD ["/authsvc", "start"]
