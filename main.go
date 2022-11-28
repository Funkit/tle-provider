package main

import (
	"github.com/Funkit/tle-provider/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
