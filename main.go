package main

import (
	"flag"
	"go-dfms-gateway/server"
	"log"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	cfgPath := flag.String("cfg", "", "Path to config file")
	address := flag.String("addr", "", "Sets gateway address")
	addressAPI := flag.String("api-addr", "", "Sets API address")
	flag.Parse()

	s := server.NewGateway(
		server.WithAddress(*address),
		server.WithAPI(*addressAPI),
		server.Debug(*debug),
		server.ConfigPath(*cfgPath),
	)

	log.Fatal(s.Start())
}
