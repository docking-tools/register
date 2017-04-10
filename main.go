package main

import (
	"flag"
	"fmt"
	config "github.com/docking-tools/register/config"
	doc "github.com/docking-tools/register/docker"
	template "github.com/docking-tools/register/template"
	"log"
	"os"

	"github.com/docking-tools/register/api"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Printf("Starting register ...")

	//    configDir := flag.String("c", config.ConfigDir(), "Path for config dir (default $DOCKING_CONFIG)")

	flag.Usage = func() {
		//fmt.Fprintf(os.Stderr, "Usage of .....", os.Args[0])
		// @TODO create Usage helper
		flag.PrintDefaults()
	}

	configFile, e := config.Load("")

	if e != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading config file:%v\n", e)
	}

	flag.StringVar(&configFile.HostIp, "ip", configFile.HostIp, "Ip for ports mapped to the host (shorthand)")

	flag.Parse()

	log.Printf("Configuration:  %v %v", len(configFile.Targets), configFile)

	clients := make([]api.RegistryAdapter, 0)
	for _, target := range configFile.Targets {
		client, err := template.NewTemplate(target)
		assert(err)
		clients = append(clients, client)
	}

	docker, err := doc.New(configFile)

	assert(err)
	docker.Start(func(status string, object interface{}, closeChan chan error) error {
		for _, client := range clients {
			if client == nil {
				continue
			}
			err := client.RunTemplate(status, object)
			if err != nil {
				log.Printf("Error on RunTemplate %v", err)
				closeChan <- err
				continue
			}
		}
		return nil
	})
}
