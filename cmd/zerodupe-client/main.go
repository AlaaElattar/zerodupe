package main

import (
	"log"
	"zerodupe/pkg/client/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}