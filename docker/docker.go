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
    config		*config.ConfigFile
    servicesMap map[string][]*api.Service	// store Key=containerId / Value=List of ServiceName
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
        config: config,
        servicesMap: make(map[string][]*api.Service),
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

	parseService := func(id string, status string) {
		container, err := doc.docker.ContainerInspect(context.Background(), id)
		if err != nil {
			closeChan <- err
			return
		}
		if (status=="") {
			status =container.State.Status
		}
	
		services :=  make([]*api.Service,0)
		if data, ok := doc.servicesMap[id]; ok {
			services =data 
		} else {
			services, err = createService(doc.config , &container)
			if err != nil {
				closeChan <- err
			}
			if services != nil {
				doc.servicesMap[id]=services
			}
		}

		for _,service := range services {
			go ep(status, service , closeChan)
		}
	
	}	

	
	// getContainerList simulates creation event for all previously existing
	// containers.
	getContainerList := func() {
		options := types.ContainerListOptions{
			Quiet: false,
		}
		cs, err := doc.docker.ContainerList(context.Background(), options)
		if err != nil {
			closeChan <- err
		}
		for _, container := range cs {
			parseService(container.ID, container.State)
		}
	}
	

	eh := eventHandler{handlers: make(map[string]func(eventtypes.Message))}
		eh.Handle("*", func(e eventtypes.Message) {
			log.Printf("new event %v id: %v",e.Status,e.ID)
			parseService(e.ID, e.Status)
		})
	
	
	
	// start listening event
	started := make(chan struct{})
	eventChan := make(chan eventtypes.Message)
	go eh.Watch(eventChan)
	go monitorContainerEvents(started, eventChan)
	defer close(eventChan)
	<-started
	
	getContainerList ()
   
   
   
   
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