FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application, creating a static binary
# CGO_ENABLED=0 is important for creating a static binary on Alpine
#RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/server ./cmd/server
RUN go build -v -o /app/server ./cmd/server

# Step 2: Use a minimal 'distroless' or 'alpine' image for the final container
# This results in a much smaller and more secure final image.
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/server /server

# Expose the port the application runs on
EXPOSE 8080

# The command to run the application
CMD [ "/server" ]