FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go app. CGO_ENABLED=0 for fully static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-jwks-validator .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/web-jwks-validator .

CMD ["./web-jwks-validator"]
