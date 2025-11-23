FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o note-go cmd/main.go

EXPOSE 50051
CMD ["./note-go"]
