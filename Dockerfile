FROM golang:1.10.3-alpine
WORKDIR /go/src/cmccacher
COPY . /go/src/cmccacher
RUN go build -o main
CMD ["/go/src/cmccacher/main"]