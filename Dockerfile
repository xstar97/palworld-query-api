# Stage 1 - Build the Go application and download CLI
FROM golang:1.22.0-alpine3.19 AS builder

# Install necessary build dependencies
RUN apk --no-cache add --update gcc musl-dev

# Create the necessary directories
RUN mkdir -p /build /output /app

# Set the working directory
WORKDIR /build

# Copy all files from the cmd directory
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY internal/routes ./internal/routes
COPY internal/config ./internal/config
COPY cmd/main.go ./main.go

# Download dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=1 go build -ldflags "-w -s" -o /output/palworld-query-api .

# Stage 2 - Create the final image
FROM alpine:3.19.1 AS runner

# Set maintainer label
LABEL maintainer="Xstar97 <dev.xstar97@gmail.com>"

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /output/palworld-query-api ./

# Create the necessary directories
RUN mkdir -p /config /logs

# Set user and group environment variables
ENV APP_USER=apps \
    APP_GROUP=apps \
    APP_USER_ID=568 \
    APP_GROUP_ID=568

# Create a non-root user and group
RUN addgroup -g $APP_GROUP_ID -S $APP_GROUP && \
    adduser -u $APP_USER_ID -S $APP_USER -G $APP_GROUP

# Change ownership of the /config directory to the non-root user and group
RUN chown -R $APP_USER:$APP_GROUP /config

# Change ownership of the /logs directory to the non-root user and group
RUN chown -R $APP_USER:$APP_GROUP /logs

# Set environment variables
ENV RCON_CLI_CONFIG=/config/rcon.yaml \
    LOGS_PATH=/logs \
    PORT=3000

# Expose the port
EXPOSE $PORT

# Set the default command to run the binary
CMD sh -c "./palworld-query-api"