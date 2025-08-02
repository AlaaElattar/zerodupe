package main

import (
	"log"
	"zerodupe/internal/server/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
