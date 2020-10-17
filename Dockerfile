FROM golang:1.15.3

WORKDIR /github.com/vadimdoga/Distributed_Systems_Lab_1
COPY . .

RUN go get -d -v

RUN go install -v

ENV LIMIT=20
ENV TIMEOUT=1s
ENV GATEWAY_ADDR=https://httpbin.org/post
ENV IP=0.0.0.0
ENV PORT=5000

EXPOSE 5000

CMD go run .
