FROM golang:alpine3.18 as builder

# Time zone
RUN apk add --no-cache tzdata

# Set necessary environmet variables needed for our images
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download


# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /dist/main /

COPY ./.env /.env
ENV TZ=Asia/Jakarta
ENV MODE=PROD

# Run executable
CMD ["./main"]