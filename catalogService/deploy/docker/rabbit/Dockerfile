FROM golang:1.20 as build

WORKDIR /build

COPY Makefile .
COPY ./internal/ /build/internal/
COPY go.mod .
COPY ./docs/ /build/docs/
COPY ./cmd/ /build/cmd/

RUN make download prebuild build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/gateapp .
COPY --from=build /build/dist/linux/coreapp .

COPY config/config.yaml /config.yaml

CMD ["/gateapp", "start"]