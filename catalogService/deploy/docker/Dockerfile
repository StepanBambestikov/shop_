FROM golang:1.20 as build

WORKDIR /build

COPY catalogService/Makefile .
COPY ./catalogService/internal/ /build/internal/
COPY ./catalogService/go.mod .
COPY ./catalogService/docs/ /build/docs/
COPY ./catalogService/cmd/ /build/cmd/

RUN make download prebuild build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/gateapp .
COPY --from=build /build/dist/linux/coreapp .

COPY catalogService/config/config.yaml /config.yaml

CMD ["/gateapp", "start"]