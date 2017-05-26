package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	config "github.com/docking-tools/register/config"
	doc "github.com/docking-tools/register/docker"
	template "github.com/docking-tools/register/template"
	"os"

	"github.com/docking-tools/register/api"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Info("Starting register ...")

	//    configDir := flag.String("c", config.ConfigDir(), "Path for config dir (default $DOCKING_CONFIG)")

	flag.Usage = func() {
		//fmt.Fprintf(os.Stderr, "Usage of .....", os.Args[0])
		// @TODO create Usage helper
		flag.PrintDefaults()
	}
	showVersion := true

	configFile, e := config.Load("")

	if e != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading config file:%v\n", e)
	}
	flag.BoolVar(&showVersion, "version", false, "Print version information and quit")
	flag.StringVar(&configFile.HostIp, "ip", configFile.HostIp, "Ip for ports mapped to the host (shorthand)")

	flag.Parse()

	if showVersion {
		fmt.Printf("Register version %s, build %s\n", Version, GitCommit)
		return
	}

	log.WithFields(log.Fields{
		"NumberOfTarget": len(configFile.Targets),
		"Config":         configFile,
	}).Info("Configuration")

	clients := make([]api.RegistryAdapter, 0)
	for _, target := range configFile.Targets {
		client, err := template.NewTemplate(target)
		assert(err)
		clients = append(clients, client)
	}

	docker, err := doc.New(configFile)

	assert(err)
	if len(clients) == 0 {
		log.Fatal("No template found.")
	}

	docker.Start(func(status string, object interface{}, closeChan chan error) error {
		for _, client := range clients {
			if client == nil {
				continue
			}
			err := client.RunTemplate(status, object)
			if err != nil {
				log.Warn("Error on RunTemplate %v", err)
				closeChan <- err
				continue
			}
		}
		return nil
	})
}
