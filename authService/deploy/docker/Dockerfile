FROM golang:1.20 as build

WORKDIR /build

COPY ./authService/Makefile .
COPY ./authService/internal/ /build/internal/
COPY ./authService/go.mod .
COPY ./authService/docs/ /build/docs/
COPY ./authService/cmd/ /build/cmd/

RUN make download prebuild doc build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/goauth .

COPY authService/config/config.yaml /config.yaml

CMD ["/goauth", "start"]
