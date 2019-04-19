FROM golang:1.12.4

WORKDIR /go/src/target-project
COPY . .

RUN go build

CMD ["target-project"]