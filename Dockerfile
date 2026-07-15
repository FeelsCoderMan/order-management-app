FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:3.19

COPY --from=builder /src/dist/server /server
CMD ["/server"]
