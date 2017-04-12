package docker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	swarm "github.com/docker/docker/api/types/swarm"
	api "github.com/docking-tools/register/api"
	"github.com/docking-tools/register/config"
	"net"
	"strconv"
)

func createService(config *config.ConfigFile, container *types.ContainerJSON, swarmService *swarm.Service) ([]*api.Service, error) {

	ports := make(map[string]DockerServicePort)
	isSwarmMode := false
	if swarmService == nil {
		// For non swarm mode container
		if container.NetworkSettings != nil && len(container.NetworkSettings.Ports) > 0 {
			// Extract runtime port mappings, relevant when using --net=bridge
			for port, published := range container.NetworkSettings.Ports {
				ports[string(port)] = servicePort(container, port, published)
			}

		}
	} else {
		isSwarmMode = true
		// For swarm mode container
		for _, port := range swarmService.Endpoint.Ports {
			ports[string(port.TargetPort)] = serviceSwarmPort(container, swarmService, port)
		}

	}
	services := make([]*api.Service, 0)

	for _, port := range ports {
		service := newService(config, swarmService, port, len(ports) > 1, isSwarmMode)
		if service == nil {
			log.WithFields(log.Fields{
				"containerID": port.ContainerID,
				"port":        port.ExposedPort,
			}).Info("ignored service")
			continue
		}
		log.WithFields(log.Fields{
			"containerID":      service.ID,
			"serviceName":      service.Name,
			"serviceSwarmMode": service.SwarmMode,
			"serviceVersion":   service.Version,
			"servicePort":      service.Port,
		}).Debug("serviceSwarmPort")
		services = append(services, service)
	}

	return services, nil
}

func newService(config *config.ConfigFile, swarmService *swarm.Service, port DockerServicePort, isgroup bool, isSwarmMode bool) *api.Service {
	log.WithFields(log.Fields{
		"containerID": port.ContainerID,
		"name":        port.ServiceName,
		"isGroup":     isgroup,
	}).Debug("serviceSwarmPort servicePort")

	container := port.container

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

	metadata, metadataFromPort := serviceMetaData(port.ExposedPort, envArrayToMap(container.Config.Env), container.Config.Labels, swarmService.Spec.Labels)

	ignore := mapDefault(metadata, "ignore", "")
	if ignore != "" {
		return nil
	}

	service := new(api.Service)
	service.Origin = port.ServicePort
	service.ID = container.ID[:12]
	service.Name = mapDefault(metadata, "name", port.ServiceName)
	service.Version = mapDefault(metadata, "version", "default")
	service.SwarmMode = isSwarmMode
	if isgroup && !metadataFromPort["name"] {
		service.Name += "-" + port.ExposedPort
	}
	var p int

	service.IP = port.HostIP
	p, _ = strconv.Atoi(port.HostPort)

	service.Port = p

	if port.PortType == "udp" {
		service.Tags = combineTags(
			mapDefault(metadata, "tags", ""), "", "udp")
		service.ID = service.ID + ":udp"
	} else {
		service.Tags = combineTags(
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

	return service
}
