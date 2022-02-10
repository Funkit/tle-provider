package main

import (
	"fmt"

	"github.com/Funkit/tle-provider/utils"
)

func main() {
	config, err := utils.GetConfiguration()
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
