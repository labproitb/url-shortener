FROM golang:1.14.2-alpine3.11

WORKDIR /url-shortener
COPY . .

RUN go build -o main ./app

CMD ["./main"]
