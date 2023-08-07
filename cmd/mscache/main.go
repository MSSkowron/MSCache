package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/MSSkowron/MSCache/internal/cache"
	"github.com/MSSkowron/MSCache/internal/node"
)

var (
	// ErrServerListenAddressNotSpecified is returned when the server's listen addres has not been specified through the flag or environment variable.
	ErrServerListenAddressNotSpecified = errors.New("server listen address flag is empty & MSCACHE_LISTENADDR environment variable is not set")
)

func main() {
	listenAddr, leaderAddr, err := parseFlags()
	if err != nil {
		fmt.Printf("Failed to read server listen address: %s\n", err.Error())
		os.Exit(1)
	}

	cache := cache.NewInMemoryCache()

	if err := node.New(listenAddr, leaderAddr, leaderAddr == "", cache).Run(); err != nil {
		fmt.Printf("Failed to start node: %s\n", err.Error())
		os.Exit(1)
	}
}

func parseFlags() (listenAddr, leaderAddr string, err error) {
	flag.StringVar(&listenAddr, "listenaddr", os.Getenv("MSCACHE_LISTENADDRESS"), "listen address of the server")
	flag.StringVar(&leaderAddr, "leaderaddr", os.Getenv("MSCACHE_LEADERADDRESS"), "listen address of the leader server")
	flag.Parse()

	if listenAddr == "" {
		err = ErrServerListenAddressNotSpecified
	}

	return
}
