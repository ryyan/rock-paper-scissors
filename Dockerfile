FROM golang:alpine

EXPOSE 5000
EXPOSE 5001

WORKDIR /app
COPY . .

# Install Node
RUN apk update \
    && apk add --no-cache git nodejs npm

# Build client
RUN cd client \
    && npm i

# Build server
RUN cd server \
    && go get -d \
    && go build -buildvcs=false

CMD ["./server/rock-paper-scissors"]
