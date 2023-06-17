FROM golang:1.20 as build

WORKDIR /build

COPY authService/Makefile .
COPY ./authService/internal/ authService/build/internal/
COPY authService/go.mod .
COPY ./authService/docs/ authService/build/docs/
COPY ./authService/cmd/ authService/build/cmd/

RUN make download prebuild build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/authsvc .

COPY authService/config/config.yaml authService/config.yaml

CMD ["/authsvc", "start"]
