package main

import (
	"flag"
	"log"

	"github.com/TheCreeper/wuu/wuu"
)

// flags
var (
	dbname string
	addr   string
)

func init() {
	flag.StringVar(&dbname, "db", "", "filepath of the database directory")
	flag.StringVar(&addr,
		"addr",
		":8080",
		"address of the interface to listen on")
	flag.Parse()
}

func main() {
	if len(dbname) == 0 {
		println("wuu: missing database filepath")
		flag.PrintDefaults()
		return
	}

	if len(addr) == 0 {
		println("wuu: no address specified")
		flag.PrintDefaults()
		return
	}

	if err := wuu.Listen(dbname, addr); err != nil {
		log.Fatal(err)
	}
}
