# Use an official Go runtime as a parent image
FROM golang:1.22

# Create a non-root user
RUN useradd -m -s /bin/bash gouser

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies as root
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o gonvm .

# Install gonvm to a directory in PATH
RUN mv gonvm /usr/local/bin/

# Change ownership of the /app directory and Go cache directories to gouser
RUN chown -R gouser:gouser /app /go /usr/local/bin/gonvm

# Switch to the new user
USER gouser

# Set GOPATH for the non-root user
ENV GOPATH /home/gouser/go
ENV PATH $GOPATH/bin:$PATH

# Run tests
CMD ["go", "test", "-v", "./..."]