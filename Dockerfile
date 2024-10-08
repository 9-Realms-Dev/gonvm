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

# Set up environment variables
ENV GO_NVM_DIR="/home/gouser/.go_nvm"
ENV NVM_CURRENT="${GO_NVM_DIR}/current/bin"
ENV PATH="${NVM_CURRENT}:${PATH}"

# Add environment variables to .bashrc
RUN echo "export GO_NVM_DIR=${GO_NVM_DIR}" >> /home/gouser/.bashrc && \
    echo "export NVM_CURRENT=${NVM_CURRENT}" >> /home/gouser/.bashrc && \
    echo "export PATH=${NVM_CURRENT}:\$PATH" >> /home/gouser/.bashrc

# Set the working directory to the user's home
WORKDIR /home/gouser

# Set the default command to run an interactive shell
CMD ["/bin/bash"]