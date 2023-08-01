package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/MSSkowron/MSCache/client"
)

func main() {
	var (
		serverAddrFlag = flag.String("serveraddr", "", "listen address of the server")
	)
	flag.Parse()

	if len(*serverAddrFlag) == 0 {
		log.Fatalln("server address is empty")
	}

	c, err := client.New(*serverAddrFlag)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	for i := 0; i < 10; i++ {
		var (
			key   = []byte(fmt.Sprintf("key:%d", i))
			value = []byte(fmt.Sprintf("value:%d", i))
		)

		if err := c.Set(context.TODO(), key, value, 100); err != nil {
			log.Fatalln(err)
		}

		time.Sleep(1 * time.Second)

		// err := c.Delete(context.TODO(), key)
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// time.Sleep(1 * time.Second)

		val, err := c.Get(context.TODO(), key)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("GOT: %s\n", string(val))
	}
}
