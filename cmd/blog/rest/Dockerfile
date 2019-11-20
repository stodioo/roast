FROM golang:1.12 as builder

WORKDIR /github.com/stodioo/roast/

COPY go.mod go.sum ./
RUN go mod download

ARG revisionID=unknown
ARG buildTimestamp=unknown

# Copy all resources
COPY . .

# Build it
RUN CGO_ENABLED=0 go build -o ./build/package/blogrest \
    -ldflags="-s -w -X main.revisionID=${revisionID} -X main.buildTimestamp=${buildTimestamp}" \
    ./cmd/blog/rest

# Build the runtime image
FROM alpine:3.10
RUN apk add --no-cache ca-certificates
WORKDIR /root

COPY --from=builder /github.com/stodioo/roast/build/package/blogrest ./service

EXPOSE 8080

ENTRYPOINT ["./service"]