package app

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

func Run() error {
	listenAddr, leaderAddr, err := parseFlags()
	if err != nil {
		return fmt.Errorf("failed to read server listen address: %s\n", err.Error())
	}

	cache := cache.NewInMemoryCache()

	if err := node.New(listenAddr, leaderAddr, leaderAddr == "", cache).Run(); err != nil {
		return fmt.Errorf("failed to start node: %s\n", err.Error())
	}

	return nil
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
