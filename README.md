# mscache

## About the project

mscache is a distributed key-value cache written in Go. Nodes communicate with each other using TCP protocol.

## How to use

In order to run a new leader node we can use _runleader_ build in _Makefile_. \
Using flag --listenaddr we can specify an address on which the node will be listening to incoming commands. Port 3000 is the default one. We can change it if we want. \
Use this command to start the leader node using _Makefile_:

```
make runleader
```

Use this command to start the leader node without using _Makefile_:

```
go build -o bin/mscache
./bin/mscache --listenaddr :3000
```

In order to run a new follower node we can use _runfollower_ build in _Makefile_. \
Using flag --listenaddr we can specify an addres on which the node will be listening to commands. Port 4000 is the default one. We can change it if we want. \
Using flag --leaderaddr we can specify a leader's address. Port 3000 is the default one.\
Use this command to start the leader node using _Makefile_:

```
make runfollower
```

Use this command to start the leader node without using _Makefile_:

```
go build -o bin/mscache
./bin/mscache --listenaddr :4000 --leaderaddr :3000
```

Each node needs to have a different listen address.
