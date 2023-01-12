package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MSSkowron/mscache/cache"
	"github.com/MSSkowron/mscache/server"
)

func main() {
	go func() {
		// Send SET
		time.Sleep(time.Second * 2)

		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.Write([]byte("SET Foo Bar 2500000000000"))
		if err != nil {
			log.Fatalln(err)
		}

		// Send GET
		time.Sleep(time.Second * 2)

		_, err = conn.Write([]byte("GET Foo"))
		if err != nil {
			log.Fatalln(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(buf[:n]))
	}()

	server := server.NewServer(":3000", true, cache.NewCache())
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
