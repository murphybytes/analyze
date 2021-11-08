package main

import (
	"flag"
	"log"
	"time"


	"github.com/osquery/osquery-go"
)

func main() {
	var osQuerySockPath, configPath string

	flag.StringVar(&configPath, "config", "conf.yaml", "the path to the yaml configuration defining checks")
	flag.StringVar(&osQuerySockPath, "osquery",  "/var/osquery/osquery.em", "the path to the osquery socket")
	flag.Parse()

	client, err := osquery.NewClient(osQuerySockPath, 10*time.Second)
	if err != nil {
		log.Fatalf("could not create osquery client %q", err )
	}
	defer client.Close()




}
