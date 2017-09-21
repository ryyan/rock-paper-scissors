FROM alpine:latest

ENV GOPATH /go
EXPOSE 5000
EXPOSE 5001

WORKDIR /app
COPY . .

# Install Go and Node
RUN apk update \
    && apk add --no-cache git go gcc g++ nodejs nodejs-npm

# Build client
RUN cd client \
    && npm install \
    && npm run build

# Build server
RUN cd server \
    && go get -d \
    && go build

CMD ["./server/server"]
