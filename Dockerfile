# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go vet ./... && go test -v ./... && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/hexclock .

FROM gcr.io/distroless/static-debian12:nonroot
LABEL org.opencontainers.image.source="https://github.com/fregas/agentbench-hexclock-v3-17-b" \
      org.opencontainers.image.description="hexclock: current UTC time as JSON"
COPY --from=build /out/hexclock /hexclock
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/hexclock"]
