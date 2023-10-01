FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o main src/main.go

EXPOSE PORT

CMD ["./main"]