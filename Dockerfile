# Using base image of Golang
FROM golang:1.24

# Set destination for COPY
WORKDIR /app

COPY go.mod ./
RUN go mod download

# Copy the source code to container
COPY . .

# Install dependency
RUN go mod tidy

# Build binary
RUN go build -o main .
EXPOSE 8080

# เรียก binary เมื่อ container start
CMD ["./main"]

