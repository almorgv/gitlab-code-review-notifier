FROM golang:1.14 as build

ENV CGO_ENABLED 0

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main ./cmd

FROM scratch
COPY /internal/database/migrations /internal/database/migrations
COPY /pkg/notifier/templates /pkg/notifier/templates
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /build/main /
EXPOSE 8080
ENTRYPOINT ["/main"]
