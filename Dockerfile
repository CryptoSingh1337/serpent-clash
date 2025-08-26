# Stage 1: Build Vue Application
FROM node:22-alpine AS frontend-builder

WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY client/package*.json ./

# Install dependencies for Vue
RUN npm install

# Copy the rest of the Vue app
COPY client/ ./

ARG ARG_VITE_DEBUG_MODE=false

ENV VITE_DEBUG_MODE=$ARG_VITE_DEBUG_MODE

# Build the Vue app for production
RUN npm run build

# Stage 2: Build Go Application
FROM golang:1.25-alpine AS backend-builder

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
RUN GOEXPERIMENT=greenteagc,jsonv2 go build -o ./server/main ./cmd

# Stage 3: Final Image
FROM alpine:latest

WORKDIR /app

# Copy the built Go binary
COPY --from=backend-builder /app/server ./server

# Copy the Vue dist folder
COPY --from=backend-builder /app/client/dist ./client/dist

ENV GO_ENV=PROD
ENV SERVER_ADDR=0.0.0.0
ENV SERVER_PORT=8080
ENV DIST_DIR=/app/client/dist
ENV DEBUG_MODE=false

# Expose the port the Go server listens on
EXPOSE 8080

# Command to run the Go server
CMD ["./server/main"]
