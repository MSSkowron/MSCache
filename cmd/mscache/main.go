package main

import (
	"github.com/MSSkowron/MSCache/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
