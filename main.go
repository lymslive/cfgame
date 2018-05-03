package main

import (
	"fmt"
	"log"

	"github.com/lymslive/cfgame/cmdline"
	"github.com/lymslive/cfgame/server"
)

func main() {
	cfg := cmdline.ParseConfig()

	var address = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("will Serve on:", address)

	log.SetPrefix("[CFGAME] ")
	log.SetFlags(log.Ltime | log.Lshortfile)

	server.Start()

	log.Printf("Serve on %s success! But cannot reach here\n", address)
}
