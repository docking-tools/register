package docker

import (
	"github.com/docker/docker/api/types"
	eventtypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	api "github.com/docking-tools/register/api"
	"github.com/docking-tools/register/config"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"time"
)

type DockerRegistry struct {
	docker       *client.Client
	events       <-chan *io.ReadCloser
	config       *config.ConfigFile
	servicesMap  map[string][]*api.Service // store Key=containerId / Value=List of ServiceName
	graphMetaMap map[string]api.Recmap     //  store Key=containerId / Value=List of Graph
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
	log.Printf("Start docker client with env %s", os.Getenv("DOCKER_HOST"))

	// Init client docker

	docker, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}

	return &DockerRegistry{
		docker:       docker,
		events:       nil,
		config:       config,
		servicesMap:  make(map[string][]*api.Service),
		graphMetaMap: make(map[string]api.Recmap),
	}, nil
}

func (doc *DockerRegistry) Start(ep api.EventProcessor) {

	// check if an error occur during watch event
	closeChan := make(chan error)

	monitorContainerEvents := func(started chan<- struct{}, c chan<- eventtypes.Message) {
		f := filters.NewArgs()
		f.Add("type", "container")
		options := types.EventsOptions{
			Filters: f,
		}
		events, errs := doc.docker.Events(context.Background(), options)
		// Whether we successfully subscribed to events or not, we can now
		// unblock the main goroutine.
		close(started)
		log.Printf("Start listenig docker event")

		for {
			select {
			case event := <-events:
				c <- event
				return
			case err := <-errs:
				if err == io.EOF {
					return
				}
				closeChan <- err
				return
			}
		}
	}

	parseService := func(config *config.ConfigFile, container *types.ContainerJSON, status string) []*api.Service {

		services := make([]*api.Service, 0)
		services, err := createService(config, container)
		if err != nil {
			closeChan <- err
		}
		if data, ok := doc.servicesMap[container.ID]; ok && len(services) == 0 {
			services = data
			delete(doc.servicesMap, container.ID)
		} else {
			doc.servicesMap[container.ID] = services
		}

		return services

	}

	parseHierarchicalMetadata := func(config *config.ConfigFile, container *types.ContainerJSON, status string) api.Recmap {

		graph := make(api.Recmap)
		graph = graphMetaData(container.Config)
		if data, ok := doc.graphMetaMap[container.ID]; ok && len(graph) == 0 {
			graph = data
			delete(doc.graphMetaMap, container.ID)
		} else {
			doc.graphMetaMap[container.ID] = graph
		}
		return graph

	}

	parseContainer := func(id string, status string) {
		container, err := doc.docker.ContainerInspect(context.Background(), id)
		if err != nil {
			closeChan <- err
			return
		}
		if status == "" {
			status = container.State.Status
		}
		instance := api.Instance{
			Services:      parseService(doc.config, &container, status),
			MetaDataGraph: parseHierarchicalMetadata(doc.config, &container, status),
			Container:     container, Status: status,
		}
		go ep(status, instance, closeChan)
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
			parseContainer(container.ID, container.State)
		}
	}

	eh := eventHandler{handlers: make(map[string]func(eventtypes.Message))}
	eh.Handle("*", func(e eventtypes.Message) {
		log.Printf("new event %v id: %v", e.Status, e.ID)
		parseContainer(e.ID, e.Status)
	})

	// start listening event
	started := make(chan struct{})
	eventChan := make(chan eventtypes.Message)
	go eh.Watch(eventChan)
	go monitorContainerEvents(started, eventChan)
	defer close(eventChan)
	<-started

	getContainerList()

	for range time.Tick(500 * time.Millisecond) {
		select {
		case err, ok := <-closeChan:
			if ok {
				if err != nil {
					// this is suppressing "unexpected EOF" in the cli when the
					// daemon restarts so it shutdowns cleanly
					if err != io.ErrUnexpectedEOF {
						closeChan <- err
						log.Printf("Error on run", err)
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
