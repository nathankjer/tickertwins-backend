FROM golang:1.19.3-alpine

WORKDIR /app

COPY . /app

RUN go build

EXPOSE 8000

CMD ["./tickertwins-backend"]