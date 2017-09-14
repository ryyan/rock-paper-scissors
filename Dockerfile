FROM alpine:latest

ENV GOPATH /go
EXPOSE 5000
EXPOSE 5001

WORKDIR /app
COPY . .

RUN apk update \
    && apk add --no-cache git go gcc g++ nodejs nodejs-npm

RUN cd client \
    && npm install \
    && npm run build

RUN cd server \
    && go get -d \
    && go build

CMD ["./server/server"]
