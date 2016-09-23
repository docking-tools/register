package main

import(

    "flag"
    "fmt"
    "log"
    "os"
    doc "github.com/docking-tools/register/docker"  
    template "github.com/docking-tools/register/template"
    config "github.com/docking-tools/register/config"

)



func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main () {
   
   log.Printf("Starting register ...")
  

	
//    configDir := flag.String("c", config.ConfigDir(), "Path for config dir (default $DOCKING_CONFIG)")

   flag.Usage= func () {
       //fmt.Fprintf(os.Stderr, "Usage of .....", os.Args[0])
       // @TODO create Usage helper
       flag.PrintDefaults()
   }
   
   	configFile, e := config.Load("")

	if e != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading config file:%v\n", e)
	}   
	
	

    //flag.StringVar(&configFile.HostIp, "-hostip", configFile.HostIp,  "Ip for ports mapped to the host")
    flag.StringVar(&configFile.HostIp, "ip", configFile.HostIp,  "Ip for ports mapped to the host (shorthand)")
    //flag.StringVar(&configFile.RegisterUrl, "-register", configFile.RegisterUrl, "URL for discovery")
    flag.StringVar(&configFile.RegisterUrl, "r", configFile.RegisterUrl, "URL for discovery (shorthand)")
    //flag.StringVar(&configFile.DockerUrl, "-docker", configFile.DockerUrl, "URL for docker")
    flag.StringVar(&configFile.DockerUrl, "d", configFile.DockerUrl, "URL for docker (shorthand)")
    
    

    flag.Parse()

    log.Printf("Configuration:  ", configFile)
 
   
   client, err:= template.NewTemplate(configFile)
   assert(err)
   docker, err:= doc.New(configFile) 
   
   assert(err)

   docker.Start(func(status string, object interface{}, closeChan chan error) error {
			err := client.RunTemplate(status, object)
			if err != nil {
				closeChan <- err
				return nil
			}
			return nil
		})	
}

