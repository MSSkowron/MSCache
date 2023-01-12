package main

import (
	"log"
	"net"
	"time"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)

		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.Write([]byte("SET Foo Bar 2500"))
		if err != nil {
			log.Fatalln(err)
		}
	}()

	server := server.NewServer(opts, cache.NewCache())
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}

}
