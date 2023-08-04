# MSCache - Distributed Key-Value Cache

## About the project

MSCache is a distributed key-value cache implemented in Go, enabling a distributed caching system where server nodes communicate with each other using the TCP protocol. The client can connect to any node in the cluster using TCP for interactions.

## Technologies

- Go 1.20

## Requirements

Make sure you have Go installed on your system.

## Installation

Clone the repository:

```
git clone <https://github.com/MSSkowron/MSCache>
```

## How to run

Navigate to the project directory:

```
cd MSCache
```

### Starting a Leader Node

To launch a leader node, use the following command:
    
```
go run ./cmd/mscache/main.go --listenaddr <address>
```

For example:

```
go run ./cmd/mscache/main.go --listenaddr 127.0.0.1:5000
```

### Starting a Follower Node

To launch a follower node, use the following command:

```
go run ./cmd/mscache/main.go --listenaddr <address> --leaderaddr <address>
```

For example:

```
go run ./cmd/mscache/main.go --listenaddr 127.0.0.1:5001 --leaderaddr 127.0.0.1:5000
```

**Note**: Each node must have a unique listen address.

## Install & Run with `go install`

You can install the app with the `go install` to have it available globally:

```
go install github.com/MSSkowron/MSCache/cmd/mscache@latest
```

Then, you have the `mscache` command available globally, and you can run it by executing:

```
mscache --listenaddr <address> --leaderaddr <address>
```

## How to interact with Server

To interact with the server, you can utilize the Client structure defined in /client/client.go. This structure provides necessary methods to communicate with the cache server. Keep in mind that you can use the `Set`, `Delete`, and `GET` methods only with the leader server, while only the `GET` method is currently available in follower servers.

An example client's code is available [**here**](./examples/client/main.go). You can run it specifying the server node's address with the `serveraddr` flag.

```
go run ./examples/client/main.go --serveraddr <address>
```

For example:

```
go run ./examples/client/main.go --serveraddr 127.0.0.1:5000
```

## License

The project is available as open source under the terms of the MIT License.
