package main

import(

    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    dockerapi "github.com/fsouza/go-dockerclient"
)

var hostIp = flag.String("ip", "", "Ip for ports mapped to the host")

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main () {
   
   log.Printf("Starting register ...")
   
   flag.Usage= func () {
       fmt.Fprintf(os.Stderr, "Usage of .....", os.Args[0])
       // @TODO create Usage helper
       flag.PrintDefaults()
   }
   
   flag.Parse()
   
   if flag.NArg() != 1 {
        if flag.NArg() == 0 {
            fmt.Fprintln(os.Stderr, "Missiong required argument for registry URI. \n\n")
        } else {
            fmt.Fprintln(os.Stderr, " ", strings.Join(flag.Args()[1:], " "))
            fmt.Fprint(os.Stderr, "Options should come before the registry URI argument. \n\n")
        }
        flag.Usage()
        os.Exit(2)
   }
   
   if *hostIp != "" {
       log.Println("Forcing host IP to ", *hostIp)
   }
   
   dockerHost:= os.Getenv("DOCKER_HOST")
   if dockerHost == "" {
        os.Setenv("DOCKER_HOST", "unix:///tmp/docker.sock")
   }
   docker, err := dockerapi.NewClientFromEnv()
   assert(err)
   
   events:= make(chan *dockerapi.APIEvents)
   assert(docker.AddEventListener(events))
   log.Println("Listening for Docker events ...")
   
   quit := make(chan struct{})
   
   // Process Docker events
   for msg:= range events {
       log.Printf("New event received %v", msg)
       switch msg.Status {
           default :
                
           case "start":
                log.Printf("Start id: %v from: %v", msg.ID, msg.From)
       }
   }
   
   close(quit)
   log.Fatal("Docker event loop closed")
}