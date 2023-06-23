FROM golang:1.20 as gomodUpdateOut

WORKDIR /build

COPY ./catalogService/Makefile .
COPY ./catalogService/internal/ /build/internal/
COPY ./catalogService/go.mod .
COPY ./catalogService/go.sum .
COPY ./catalogService/docs/ /build/docs/
COPY ./catalogService/cmd/ /build/cmd/

RUN go mod tidy -go=1.20

FROM golang:1.20 as build

WORKDIR /build

COPY ./catalogService/Makefile .
COPY ./catalogService/internal/ /build/internal/
COPY --from=gomodUpdateOut ./build/go.mod .
COPY --from=gomodUpdateOut ./build/go.sum .
COPY ./catalogService/docs/ /build/docs/
COPY ./catalogService/cmd/ /build/cmd/

RUN make prebuild doc build-linux

FROM alpine:3.17

WORKDIR /

ENV GIN_MODE=release

COPY --from=build /build/dist/linux/gateapp .
COPY --from=build /build/dist/linux/coreapp .
COPY --from=build /build/go.mod .
COPY --from=build /build/go.sum .

COPY catalogService/config/config.yaml /config.yaml

CMD ["/gateapp", "start"]
