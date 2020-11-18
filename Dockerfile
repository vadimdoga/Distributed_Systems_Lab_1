FROM golang:1.15.3

WORKDIR /products_microservice

RUN go get -d -v

RUN go install -v

COPY ./products_microservice .

CMD go run .
