package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/client"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	var (
		listenAddrFlag = flag.String("listenaddr", "", "listen address of the server")
		leaderAddrFlag = flag.String("leaderaddr", "", "listen address of the leader server")
	)
	flag.Parse()

	go func() {
		time.Sleep(time.Second * 2)
		c, err := client.New(":3000")
		if err != nil {
			log.Fatalln(err)
		}

		c.Set(context.TODO(), []byte("mateusz"), []byte("skowron"), 100000000000)

		c.Close()
	}()

	if len(*listenAddrFlag) == 0 {
		panic("listenaddr flag is required")
	}

	if err := server.New(*listenAddrFlag, *leaderAddrFlag, len(*leaderAddrFlag) == 0, cache.New()).Run(); err != nil {
		panic(err)
	}
}
