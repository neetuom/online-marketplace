# Start from golang base image
FROM golang:alpine as builder

# Install git.
# Git is required for fetching the dependencies.
# RUN apk update && apk upgrade && \
#    apk add --no-cache git

# Add Maintainer Info
LABEL maintainer="Om Prakash <op_ved@hotmail.com>"


RUN mkdir /app
WORKDIR /app


# Build Args
ARG LOG_DIR=/app/logs

# Create Log Directory
RUN mkdir -p ${LOG_DIR}

# Environment Variables
ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log 


# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o customer-service



######## Start a new stage from scratch #######

# Run container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

#RUN mkdir /app
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/customer-service .

# Expose port 60060 to the outside world
EXPOSE 60061

CMD ["./customer-service"]