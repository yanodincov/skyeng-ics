# First stage: build
FROM golang:1.23.2-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to install dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application, creating an executable named 'app'
RUN go build -o app ./cmd

# Second stage: run
FROM alpine:3.20.3

# Set the working directory for the runtime container
WORKDIR /app

# Copy the compiled executable from the build stage
COPY --from=build /app/app .

# Specify the command to run the application
CMD ["./app"]