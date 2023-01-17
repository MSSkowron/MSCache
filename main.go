package main

import (
	"flag"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	var (
		listenAddrFlag = flag.String("listenaddr", "", "listen address of the server")
		leaderAddrFlag = flag.String("leaderaddr", "", "listen address of the leader server")
	)
	flag.Parse()

	if len(*listenAddrFlag) == 0 {
		panic("listenaddr flag is required")
	}

	if err := server.New(*listenAddrFlag, *leaderAddrFlag, len(*leaderAddrFlag) == 0, cache.New()).Run(); err != nil {
		panic(err)
	}
}
