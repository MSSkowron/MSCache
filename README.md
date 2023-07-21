# MSCache - Distributed Key-Value Cache

## About the project

MSCache is a distributed key-value cache implemented in Go. It enables a distributed caching system where server nodes communicate with each other using the TCP protocol.

## Technologies

- Go 1.20

## Requirements
Make sure you have Go installed on your system before running the application.

## Installation
Clone the repository:

git clone https://github.com/MSSkowron/MSCache

## How to run
Navigate to the project directory:

```
cd MSCache
```

### Running the Server

- Starting a Leader Node
    To start a leader node, run the following command:
    ```
    go run ./server/cmd/main.go --listenaddr <port>
    ```

    For example:
    ```
    go run ./server/cmd/main.go --listenaddr :5000
    ```

- Starting a Follower Node
    To start a follower node, use the following command:
    ```
    go run ./server/cmd/main.go --listenaddr <port> --leaderaddr <port>
    ```

    For example:
    ```
    go run ./server/cmd/main.go --listenaddr :5001 --leaderaddr :5000
    ```

**Note**: Each node needs to have a different listen address.

## How to Use It
To interact with the server, you can utilize the Client structure defined in ./client/client.go. This structure provides the necessary methods to communicate with the cache server.