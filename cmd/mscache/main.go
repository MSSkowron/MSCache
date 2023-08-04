package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MSSkowron/MSCache/internal/cache"
	"github.com/MSSkowron/MSCache/internal/node"
)

func main() {
	listenAddr, leaderAddr, err := parseFlags()
	if err != nil {
		handleError(err)
	}

	cache := cache.NewInMemoryCache()

	err = node.New(listenAddr, leaderAddr, leaderAddr == "", cache).Run()
	if err != nil {
		handleError(err)
	}
}

func parseFlags() (listenAddr, leaderAddr string, err error) {
	flag.StringVar(&listenAddr, "listenaddr", "", "listen address of the server")
	flag.StringVar(&leaderAddr, "leaderaddr", "", "listen address of the leader server")
	flag.Parse()

	if len(listenAddr) == 0 {
		err = fmt.Errorf("server's listen address is empty. Specify it with --listenaddr flag")
	}

	return
}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
