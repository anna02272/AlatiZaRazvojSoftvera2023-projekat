FROM golang:alpine

WORKDIR ./app

# Copy the source code from the current directory to the working directory inside the container
COPY . .
COPY go.mod /app
# Build the Go app
RUN go build -o main .

EXPOSE 8000

# Expose port 8000 to the outside world

# Command to run the executable
CMD ["./main"]
