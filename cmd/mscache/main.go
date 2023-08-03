package main

import (
	"flag"
	"log"

	"github.com/MSSkowron/MSCache/internal/cache"
	"github.com/MSSkowron/MSCache/internal/node"
)

func main() {
	var (
		listenAddrFlag = flag.String("listenaddr", "", "listen address of the server")
		leaderAddrFlag = flag.String("leaderaddr", "", "listen address of the leader server")
	)
	flag.Parse()

	if len(*listenAddrFlag) == 0 {
		log.Fatalln("listen address is empty")
	}

	if err := node.New(*listenAddrFlag, *leaderAddrFlag, len(*leaderAddrFlag) == 0, cache.New()).Run(); err != nil {
		log.Fatalln(err)
	}
}
