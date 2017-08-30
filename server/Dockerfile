FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && apk add git
RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT ["go-wrapper", "run"]
