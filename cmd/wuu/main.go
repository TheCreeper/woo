package main

import (
	"flag"
	"log"

	"github.com/TheCreeper/wuu"
)

func main() {
	addr := flag.String("addr", "", "tcp network address to listen on")
	dbname := flag.String("db", "", "filepath of the database directory")
	flag.Parse()

	if len(*addr) == 0 {
		println("wuu: no address specified")
		flag.PrintDefaults()
		return
	}

	if len(*dbname) == 0 {
		println("wuu: missing database filepath")
		flag.PrintDefaults()
		return
	}

	if err := wuu.Listen(*addr, *dbname); err != nil {
		log.Fatal(err)
	}
}
