package main

import (
	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	server := server.NewServer(opts, cache.NewCache())
	server.Run()
}
