FROM golang:1.22.1-alpine

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY main/*.go main/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./PirateAntiCyrillicBot ./main/main.go

# Run
CMD ["./PirateAntiCyrillicBot"]