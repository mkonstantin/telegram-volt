FROM golang:1.20.4

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o telegram-api telegram-api/cmd/main

RUN chmod +x telegram-api

CMD ["./telegram-api"]