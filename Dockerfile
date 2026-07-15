FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/server ./cmd

FROM alpine:3.19

COPY --from=builder /app/server /server
CMD ["/server"]
