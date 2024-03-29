# Rock, Paper, Scissors

Persistent RPS game written in React and Go with websockets

https://github.com/ryyan/rock-paper-scissors/assets/4228816/1874a17b-be3f-4ab9-8d7b-ec2f3d897749

## Setup

### Build

```
docker build --rm -t rps .
```

### Run (dev)

```
docker run -it -p 5000:5000 -p 5001:5001 -v ${PWD}:/app rps sh
cd client
npm run watch &
cd ..
./server/rock-paper-scissors
```

### Run (prod)

- Update the key/value pairs in client/brunch-config.json under the replacer plugin

```
docker run -d -p 5000:5000 -p 5001:5001 -e NODE_ENV=production rps
```
