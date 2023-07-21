package main

import (
	"flag"
	"log"

	"github.com/MSSkowron/MSCache/server"
	"github.com/MSSkowron/MSCache/server/cache"
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

	if err := server.New(*listenAddrFlag, *leaderAddrFlag, len(*leaderAddrFlag) == 0, cache.New()).Run(); err != nil {
		log.Fatalln(err)
	}
}
