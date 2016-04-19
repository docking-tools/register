package docker

import (
    "net"
    "path"
    "log"    
    "strconv"
    "strings"
    "github.com/docker/engine-api/types"
    api "github.com/docking-tools/register/api"
    "github.com/docking-tools/register/config"
)


func createService(config *config.ConfigFile, container *types.ContainerJSON) ([]*api.Service, error) {
 
	ports := make(map[string]DockerServicePort)

	if container.NetworkSettings != nil && len(container.NetworkSettings.Ports)>0 {
		// Extract runtime port mappings, relevant when using --net=bridge
		for port, published := range container.NetworkSettings.Ports {
			ports[string(port)] = servicePort(container, port, published)
		}

	} else {
		ports["die"] = DockerServicePort{
			ContainerID:       container.ID,
			ContainerHostname: container.Config.Hostname,
			container:         container,
		}
	}
	
    services := make([]*api.Service,0)

	for _, port := range ports {
        service := newService(config, port, len(ports) > 1)
		if service == nil {
				log.Println("ignored:", container.ID[:12], "service on port", port.ExposedPort)
			continue
		}
		services = append(services, service)
    }

	return services, nil
}


func  newService(config *config.ConfigFile, port DockerServicePort, isgroup bool) *api.Service {
	container := port.container
	defaultName := strings.Split(path.Base(container.Config.Image), ":")[0]

	// not sure about this logic. kind of want to remove it.
	if config.HostIp != "" {
		port.HostIP = config.HostIp
	}

	hostname := Hostname
	if hostname == "" {
		hostname = port.HostIP
	}
	if port.HostIP == "0.0.0.0" || port.HostIP == "" {
		ip, err := net.ResolveIPAddr("ip", hostname)
		if err == nil {
			port.HostIP = ip.String()
		}
	}



	metadata, metadataFromPort := serviceMetaData(container.Config, port.ExposedPort)

	ignore := mapDefault(metadata, "ignore", "")
	if ignore != "" {
		return nil
	}

	service := new(api.Service)
	service.Origin = port.ServicePort
	service.ID = container.ID[:12]
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
