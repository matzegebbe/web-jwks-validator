FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go app. CGO_ENABLED=0 for fully static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-jwks-validator .

FROM alpine:3.23

RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN adduser -D -u 1000 appuser

WORKDIR /app/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/web-jwks-validator .

# Run as non-root user
USER appuser

CMD ["./web-jwks-validator"]
