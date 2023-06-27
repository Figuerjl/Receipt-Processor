# Use the official Golang image as the base image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Fetch the code from the GitHub repository
RUN git clone https://github.com/Figuerjl/Receipt-Processor .

# Install the mux package
RUN go get -u github.com/gorilla/mux

# Install the dependencies
RUN go mod download

# Build the application
RUN go build -o main

# Set the command to run the application
CMD ["./main"]
