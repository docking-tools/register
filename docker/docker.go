package docker

import (
    "log"
    "os"
    "net"
    "path"
    "strconv"
    "strings"
    dockerapi "github.com/fsouza/go-dockerclient"
    api "github.com/docking-tools/register/api"
    "github.com/docking-tools/register/config"
)

type DockerRegistry struct {
    docker      *dockerapi.Client
    events      <-chan *dockerapi.APIEvents
    registry    api.RegistryAdapter
}

type DockerServicePort struct {
    api.ServicePort
    ContainerHostname string
	ContainerID       string
	ContainerName     string
	container         *dockerapi.Container
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(registry api.RegistryAdapter, config *config.ConfigFile) (*DockerRegistry, error) {
    
    // Init docker
    
   dockerHost:= config.DockerUrl
   if dockerHost == "" {
   	dockerHost= os.Getenv("DOCKER_HOST")
   }
   if dockerHost == "" {
        os.Setenv("DOCKER_HOST", "unix:///tmp/docker.sock")
   }
   docker, err := dockerapi.NewClientFromEnv()
   assert(err)
   
   events:= make(chan *dockerapi.APIEvents)
   assert(docker.AddEventListener(events))
   
    return &DockerRegistry{
        docker:   docker,
        events: events,
        registry: registry,
    }, nil
}

func (doc * DockerRegistry) Start() {
   log.Println("Listening for Docker events ...")
   
   quit := make(chan struct{})
   
   // Process Docker events
   for msg:= range doc.events {
       log.Printf("New event received %v", msg)

  
       log.Printf("%v id: %v from: %v", msg.Status, msg.ID, msg.From)
       doc.createService(msg.Status, msg.ID)
    }
      close(quit)
   log.Fatal("Docker event loop closed")
}



func (doc *DockerRegistry) createService(status string, containerId string) {
   container, err:= doc.docker.InspectContainer(containerId)
   	if err != nil {
		log.Println("unable to inspect container:", containerId[:12], err)
		return
	}     
	ports := make(map[string]DockerServicePort)
	
	// Extract configured host port mappings, relevant when using --net=host
	for port, published := range container.HostConfig.PortBindings {
		ports[string(port)] = servicePort(container, port, published)
	}

	// Extract runtime port mappings, relevant when using --net=bridge
	for port, published := range container.NetworkSettings.Ports {
		ports[string(port)] = servicePort(container, port, published)
	}

	if len(ports) == 0 {
		log.Println("ignored:", container.ID[:12], "no published ports")
		return
	}

	for _, port := range ports {
        service := doc.newService(port, len(ports) > 1)
		if service == nil {
				log.Println("ignored:", container.ID[:12], "service on port", port.ExposedPort)
			continue
		}
		err := doc.registry.RunTemplate(strings.ToUpper(status), service)
		if err != nil {
			log.Println("RunTemplate failed:", service, err)
			continue
		}
//		doc.services[container.ID] = append(b.services[container.ID], service)
    }
}


func (doc *DockerRegistry) newService(port DockerServicePort, isgroup bool) *api.Service {
	container := port.container
	defaultName := strings.Split(path.Base(container.Config.Image), ":")[0]

	// not sure about this logic. kind of want to remove it.
	hostname := Hostname
	if hostname == "" {
		hostname = port.HostIP
	}
	if port.HostIP == "0.0.0.0" {
		ip, err := net.ResolveIPAddr("ip", hostname)
		if err == nil {
			port.HostIP = ip.String()
		}
	}

//	if b.config.HostIp != "" {
//		port.HostIP = doc.config.HostIp
//	}

	metadata, metadataFromPort := serviceMetaData(container.Config, port.ExposedPort)

	ignore := mapDefault(metadata, "ignore", "")
	if ignore != "" {
		return nil
	}

	service := new(api.Service)
	service.Origin = port.ServicePort
	service.ID = hostname + ":" + container.Name[1:] + ":" + port.ExposedPort
	service.Name = mapDefault(metadata, "name", defaultName)
	service.Version = mapDefault(metadata, "version", "default")
	if isgroup && !metadataFromPort["name"] {
		service.Name += "-" + port.ExposedPort
	}
	var p int
//	if doc.config.Internal == true {
//		service.IP = port.ExposedIP
//		p, _ = strconv.Atoi(port.ExposedPort)
//	} else {
		service.IP = port.HostIP
		p, _ = strconv.Atoi(port.HostPort)
//	}
	service.Port = p

	if port.PortType == "udp" {
        service.Tags = combineTags(
//			mapDefault(metadata, "tags", ""), b.config.ForceTags, "udp")
			mapDefault(metadata, "tags", ""), "", "udp")
		service.ID = service.ID + ":udp"
	} else {
		service.Tags = combineTags(
//			mapDefault(metadata, "tags", ""), b.config.ForceTags)
			mapDefault(metadata, "tags", ""), "")
	}

	id := mapDefault(metadata, "id", "")
	if id != "" {
		service.ID = id
	}

	delete(metadata, "id")
	delete(metadata, "tags")
	delete(metadata, "name")
	delete(metadata, "version")
	service.Attrs = metadata
//	service.TTL = doc.config.RefreshTtl

	return service
}

var Hostname string

func init() {
	// It's ok for Hostname to ultimately be an empty string
	// An empty string will fall back to trying to make a best guess
	Hostname, _ = os.Hostname()
}