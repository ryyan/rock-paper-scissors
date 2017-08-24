# Rock, Paper, Scissors

Persistent RPS game written in React and Go with websockets

## Setup

* Install go, node, npm

* Build client

```
cd client
npm install
npm run build
```

* Build and run server

```
cd server
go get && go build
./server -port=:5000
```
