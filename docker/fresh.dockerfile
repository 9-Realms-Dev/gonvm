# Use an official Go runtime as a parent image
FROM golang:1.22

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN go mod download

# Build the application
RUN go build -o gonvm .

# Add the application to PATH
ENV PATH="/app:${PATH}"

# Create a non-root user
RUN useradd -m -s /bin/bash gouser

# Switch to the new user
USER gouser

# Set the working directory to the user's home
WORKDIR /home/gouser

# Set the default command to run an interactive shell
CMD ["/bin/bash"]