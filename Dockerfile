FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o hecksweeper-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /app/hecksweeper-server /usr/local/bin/
RUN mkdir -p /data
EXPOSE 2222
CMD ["hecksweeper-server", "--port", "2222", "--key", "/data/host_key", "--data", "/data/hecksweeper_data.json"]
