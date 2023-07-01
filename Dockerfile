FROM golang:1.20

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

EXPOSE 9900

CMD ["./main"]