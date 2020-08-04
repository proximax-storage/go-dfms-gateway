package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	gateway "github.com/proximax-storage/go-dfms-gateway/server"
	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Print("Wrong number of the arguments")
		return
	}
	address := flag.Arg(0)

	debug := flag.Bool("debug", false, "Enable debug mode")
	cfgPath := flag.String("cfg", "", "Path to config file")

	g := gateway.NewGateway(
		apihttp.NewClientAPI(address),
		gateway.Debug(*debug),
		gateway.ConfigPath(*cfgPath),
	)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		err := g.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(g.Start())
}
