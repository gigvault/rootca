FROM golang:1.23-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/rootca ./cmd/rootca

FROM alpine:3.18
RUN apk add --no-cache ca-certificates bash
COPY --from=builder /out/rootca /usr/local/bin/rootca
COPY scripts/ /scripts/
ENTRYPOINT ["/usr/local/bin/rootca"]

