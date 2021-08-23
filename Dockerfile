FROM golang:1.17.0 AS builder

RUN mkdir -p /trisongulate
WORKDIR /trisongulate

ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM ubuntu:latest

ENTRYPOINT [ "/trisongulate/trisongulate" ]
WORKDIR /root/
COPY --from=builder /trisongulate/trisongulate .
CMD ["./trisongulate"]  