FROM golang:1.15.3

WORKDIR /github.com/vadimdoga/Distributed_Systems_Lab_1
COPY ./Distributed_Systems_Lab_1 .

RUN go get -d -v

RUN go install -v

CMD go run .
