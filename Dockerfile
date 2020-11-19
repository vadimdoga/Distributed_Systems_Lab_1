FROM golang:1.15.3

WORKDIR /products_service

COPY . .

RUN go get -d -v

RUN go install -v

CMD go run .
