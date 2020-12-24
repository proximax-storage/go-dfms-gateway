package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gateway "github.com/proximax-storage/go-dfms-gateway"
	apihttp "github.com/proximax-storage/go-xpx-dfms-api-http"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	cfgPath := flag.String("cfg", "", "Path to config file")

	cors := flag.Bool("cors", false, "Enable CORs")
	allowedMethods := flag.String("methods", "", "List of allowed CORs methods separated by commas.")
	allowedHeaders := flag.String("headers", "", "List of allowed CORs headers separated by commas.")
	allowedOrigins := flag.String("origins", "", "List of allowed CORs origins separated by commas.")
	flag.Parse()

	var methods []string
	if len(*allowedMethods) > 0 {
		methods = strings.Split(*allowedMethods, ",")
	}

	var headers []string
	if len(*allowedHeaders) > 0 {
		headers = strings.Split(*allowedHeaders, ",")
	}

	var origins []string
	if len(*allowedOrigins) > 0 {
		origins = strings.Split(*allowedOrigins, ",")
	}

	if len(flag.Args()) != 1 {
		log.Print("Wrong number of arguments")
		return
	}
	address := flag.Arg(0)

	g := gateway.NewGateway(
		apihttp.NewClientAPI(address),
		gateway.Debug(*debug),
		gateway.ConfigPath(*cfgPath),
		gateway.EnableCORs(*cors),
		gateway.AllowedMethods(methods...),
		gateway.AllowedHeaders(headers...),
		gateway.AllowedOrigins(origins...),
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

	err := g.Start()
	if err != nil {
		log.Fatal(err)
	}
}
