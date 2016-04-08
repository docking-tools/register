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
  

	
    configDir := flag.String("c", config.ConfigDir(), "Path for config dir (default $DOCKING_CONFIG)")

   flag.Usage= func () {
       fmt.Fprintf(os.Stderr, "Usage of .....", os.Args[0])
       // @TODO create Usage helper
       flag.PrintDefaults()
   }
   
   flag.Parse()

   	configFile, e := config.Load(*configDir)

	if e != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading config file:%v\n", e)
	}   

    hostIp := flag.String("ip", "", "Ip for ports mapped to the host")


   
    if *hostIp != "" {
       log.Println("Forcing host IP to ", *hostIp)
    }
    log.Printf("Configuration:  ", configFile)
 
   
   client, err:= template.NewTemplate(configFile)
   assert(err)
   docker, err:= doc.New(client, configFile) 
   
   assert(err)

   docker.Start()	
}

