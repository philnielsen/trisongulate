FROM golang:1.15.1

RUN mkdir -p /trisongulate
WORKDIR /trisongulate

ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

ENTRYPOINT [ "/trisongulate/trisongulate" ]
