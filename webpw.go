package main

import (
	"fmt"
	"github.com/z0rr0/gopwgen/pwgen"
)

func main() {
	fmt.Println("start")
	pw, err := pwgen.New(
		8, 2, "", "",
		false, true, false,
		false, false, false, false, false,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(pw)
}
