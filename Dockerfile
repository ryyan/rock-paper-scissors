FROM golang:alpine

EXPOSE 5000
EXPOSE 5001

WORKDIR /app
COPY . .

# Install Node
RUN apk update \
    && apk add --no-cache git nodejs nodejs-npm

# Build client
RUN cd client \
    && npm install -g yarn \
    && yarn install \
    && yarn run build

# Build server
RUN cd server \
    && go get -d \
    && go build

CMD ["./server/server"]
