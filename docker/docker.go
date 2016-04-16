package docker

import (
	"encoding/json"
	"io"
    "log"
    "os"
    "time"
    "github.com/docker/engine-api/client"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/filters"
    eventtypes "github.com/docker/engine-api/types/events"
    "golang.org/x/net/context"
    api "github.com/docking-tools/register/api"
    "github.com/docking-tools/register/config"
    
)

type DockerRegistry struct {
    docker      *client.Client
    events      <-chan *io.ReadCloser
}

type DockerServicePort struct {
    api.ServicePort
    ContainerHostname string
	ContainerID       string
	ContainerName     string
	container         *types.ContainerJSON
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(config *config.ConfigFile) (*DockerRegistry, error) {
    
    // Init docker
    
   dockerHost:= config.DockerUrl
   if dockerHost == "" {
   	dockerHost= os.Getenv("DOCKER_HOST")
   }
	log.Printf("Start docker client %s", dockerHost)

	// Init client docker
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	docker, err := client.NewClient(dockerHost, "v1.22", nil, defaultHeaders)

   if err != nil {
        panic(err)
    }


   
    return &DockerRegistry{
        docker:   docker,
        events: nil,
    }, nil
}

func (doc * DockerRegistry) Start(ep api.EventProcessor) {

   // check if an error occur during watch event
   closeChan := make(chan error)
   
	monitorContainerEvents := func(started chan<- struct{}, c chan<-eventtypes.Message) {
		f := filters.NewArgs()
		f.Add("type", "container")
		options := types.EventsOptions{
			Filters: f,
		}
		resBody, err := doc.docker.Events(context.Background(), options)
		// Whether we successfully subscribed to events or not, we can now
		// unblock the main goroutine.
		close(started)
		log.Printf("Start listenig docker event")		
		if err != nil {
			closeChan <- err
			return
		}
		defer resBody.Close()

		// Decode event
		dec := json.NewDecoder(resBody)
		for {
			var event eventtypes.Message
			err := dec.Decode(&event)
			if err != nil && err == io.EOF {
				break
			}
			c <- event
		
		}
	}
	
	
	eh := eventHandler{handlers: make(map[string]func(eventtypes.Message))}
		eh.Handle("*", func(e eventtypes.Message) {
			filters := filters.NewArgs()
			filters.Add("id", e.ID)
			container, err := doc.docker.ContainerInspect(context.Background(), e.ID)
			if err != nil {
				closeChan <- err
			}
			log.Printf("new event %v id: %v",e.Status,e.ID)
			services :=  make([]*api.Service,0)
			if &container != nil  {
				service := new(api.Service)
				service.ID=e.ID[:12]
				services = append( make([]*api.Service,0), service)
			} else {
				services, err = createService(&container)
				if err != nil {
					closeChan <- err
				}
			}

			for _,service := range services {
				go ep(e.Status, service , closeChan)
			}
		})
	
	
	
	// start listening event
	started := make(chan struct{})
	eventChan := make(chan eventtypes.Message)
	go eh.Watch(eventChan)
	go monitorContainerEvents(started, eventChan)
	defer close(eventChan)
	<-started
   
   
   
   
   for range time.Tick(500 * time.Millisecond) {
		select {
		case err, ok := <-closeChan:
			if ok {
				if err != nil {
					// this is suppressing "unexpected EOF" in the cli when the
					// daemon restarts so it shutdowns cleanly
					if err != io.ErrUnexpectedEOF {
						closeChan <- err
						log.Printf("Error on run",err)	
						close(closeChan)
					}
				}
			}
		default:
			// just skip
		}
   }
   log.Fatal("Docker event loop closed")
}




var Hostname string

func init() {
	// It's ok for Hostname to ultimately be an empty string
	// An empty string will fall back to trying to make a best guess
	Hostname, _ = os.Hostname()
}