# MSCache - Distributed Key-Value Cache

## About the project

MSCache is a distributed key-value cache implemented in Go, facilitating a distributed caching system where server nodes communicate with each other using the TCP protocol. The client can connect to any node in the cluster using TCP for interactions.

## Technologies

- Go 1.20

## Requirements

Make sure you have Go installed on your system before running the application.

## Installation

Clone the repository:

git clone <https://github.com/MSSkowron/MSCache>

## How to run

Navigate to the project directory:

```
cd MSCache
```

### Running the Server

- Starting a Leader Node
    To launch a leader node, use the following command:

    ```
    go run ./server/cmd/main.go --listenaddr <address>
    ```

    For example:

    ```
    go run ./server/cmd/main.go --listenaddr 127.0.0.1:5000
    ```

- Starting a Follower Node
    To launch a follower node, use the following command:

    ```
    go run ./server/cmd/main.go --listenaddr <address> --leaderaddr <address>
    ```

    For example:

    ```
    go run ./server/cmd/main.go --listenaddr 127.0.0.1:5001 --leaderaddr 127.0.0.1:5000
    ```

**Note**: Each node must have a unique listen address.

## How to Use It

To interact with the server, you can utilize the Client structure defined in /client/client.go. This structure provides necessary methods to communicate with the cache server. Keep in mind that you can use the `Set`, `Delete`, and `GET` methods only with the leader server, while only the `GET` method is currently available in follower servers. 

An example client's code is available [**here**](./client/runtest/main.go).

## License

The project is available as open source under the terms of the MIT License.
