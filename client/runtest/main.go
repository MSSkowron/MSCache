package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MSSkowron/mscache/client"
)

func main() {
	c, err := client.New(":3000")
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	for i := 0; i < 10; i++ {
		var (
			key   = []byte(fmt.Sprintf("key_%d", i))
			value = []byte(fmt.Sprintf("value_%d", i))
		)

		if err := c.Set(context.TODO(), key, value, 1000000000000000); err != nil {
			log.Fatalln(err)
		}

		time.Sleep(1 * time.Second)

		val, err := c.Get(context.TODO(), key)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("GOT: %s\n", string(val))
	}
}
