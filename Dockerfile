FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Copy taco-lib module (needed by replace directive in go.mod)
COPY taco-lib/ /src/taco-lib/

# Copy taco-vulndb module
COPY taco-vulndb/go.mod taco-vulndb/go.sum /src/taco-vulndb/
WORKDIR /src/taco-vulndb
RUN go mod download

COPY taco-vulndb/ /src/taco-vulndb/
RUN CGO_ENABLED=0 go build -o /bin/taco-vulndb ./cmd/taco-vulndb

FROM alpine:3.23

RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/taco-vulndb /usr/local/bin/taco-vulndb

ENTRYPOINT ["taco-vulndb"]
