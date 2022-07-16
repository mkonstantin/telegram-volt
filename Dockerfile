FROM golang:1.18.2

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./

RUN go mod download

COPY ./cmd/main/main.go .

RUN go build -o /telegram-api

CMD [ "/telegram-api" ]