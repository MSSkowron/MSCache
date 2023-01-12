package main

import (
	"flag"
	"log"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	var (
		listenAddrFlag = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddrFlag = flag.String("leaderaddr", "", "listen address of the leader server")
	)
	flag.Parse()

	if err := server.NewServer(*listenAddrFlag, *leaderAddrFlag, len(*leaderAddrFlag) == 0, cache.NewCache()).Run(); err != nil {
		log.Fatalln(err)
	}
}
