# MSCache - Distributed Key-Value Cache

## About the project

MSCache is a distributed key-value cache implemented in Go, designed to facilitate a distributed caching system where server nodes communicate with each other using the TCP protocol. The client can connect to any node within the cluster using TCP for interactions.

## Technologies

- Go 1.20

## Requirements

Ensure that you have Go installed on your system.

## Installation

Clone the repository:

```
git clone <https://github.com/MSSkowron/MSCache>
```

## Getting Started

Navigate to the project directory:

```
cd MSCache
```

### Launching a Leader Node

To launch a leader node, use the following command:
    
```
go run ./cmd/mscache/main.go --listenaddr <address>
```

For example:

```
go run ./cmd/mscache/main.go --listenaddr 127.0.0.1:5000
```

### Launching a Follower Node

To launch a follower node, use the following command:

```
go run ./cmd/mscache/main.go --listenaddr <address> --leaderaddr <address>
```

For example:

```
go run ./cmd/mscache/main.go --listenaddr 127.0.0.1:5001 --leaderaddr 127.0.0.1:5000
```

**Note**: Each node must have a unique listen address.

## Install & Run using `go install`

Install the application globally using `go install`:

```
go install github.com/MSSkowron/MSCache/cmd/mscache@latest
```

With this, the mscache command becomes accessible globally. Run it using:

```
mscache --listenaddr <address> --leaderaddr <address>
```

## How to interact with Server

To interact with the server, you can utilize the `Client` structure defined in `/client/client.go`. This structure provides methods to communicate with the cache server. `Set`, `Delete`, and `GET` methods are exclusively available with the leader server, while only the `GET` method is presently available with follower servers.

An example client's code is provided [**here**](./examples/client/main.go). You can run it specifying the server node's address with the `serveraddr` flag.

```
go run ./examples/client/main.go --serveraddr <address>
```

For example:

```
go run ./examples/client/main.go --serveraddr 127.0.0.1:5000
```

## License

The project is available as open source under the terms of the MIT License.
