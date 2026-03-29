FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/taco-vulndb ./cmd/taco-vulndb

FROM alpine:3.23

RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/taco-vulndb /usr/local/bin/taco-vulndb

ENTRYPOINT ["taco-vulndb"]
