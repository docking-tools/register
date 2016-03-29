package main

import(

    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    doc "github.com/docking-tools/register/docker"  
    template "github.com/docking-tools/register/template"  

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

 
   
   client, err:= template.NewTemplate(strings.TrimSuffix(flag.Arg(0), "/"))
   assert(err)
   docker, err:= doc.New(client) 
   
   assert(err)

   docker.Start()	
}

