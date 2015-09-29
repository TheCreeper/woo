package main

import (
	"log"

	"github.com/TheCreeper/wuu/wuu"
)

func main() {
	if err := wuu.Listen("./new.db", ":8080"); err != nil {
		log.Fatal(err)
	}
}
