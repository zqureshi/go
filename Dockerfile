FROM golang:alpine

ADD . /go/src/github.com/zqureshi/go
WORKDIR /go/src/github.com/zqureshi/go

RUN go build -o /app/go .

CMD ["/app/go"]
