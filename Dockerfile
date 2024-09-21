# Stage 1: Build Vue Application
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY client/package*.json ./

# Install dependencies for Vue
RUN npm install

# Copy the rest of the Vue app
COPY client/ ./

# Build the Vue app for production
RUN npm run build

# Stage 2: Build Go Application
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY server/go.mod server/go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the Go app
COPY server/ ./

# Copy the Vue dist folder from the frontend build
COPY --from=frontend-builder /app/dist ./client/dist

# Build the Go app
RUN go build -o ./server/main ./cmd/app

# Stage 3: Final Image
FROM alpine:latest

WORKDIR /app

# Copy the built Go binary
COPY --from=backend-builder /app/server ./server

# Copy the Vue dist folder
COPY --from=backend-builder /app/client/dist ./client/dist

ENV SERVER_ADDR=0.0.0.0
ENV SERVER_PORT=8080
ENV DIST_DIR=/app/client/dist

# Expose the port the Go server listens on
EXPOSE 8080

# Command to run the Go server
CMD ["./server/main"]
