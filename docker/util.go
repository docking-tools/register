package docker

import (
	"regexp"
	"strconv"
	"strings"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	api "github.com/docking-tools/register/api"
	"github.com/docker/go-connections/nat"
)

func mapDefault(m map[string]string, key, default_ string) string {
	v, ok := m[key]
	if !ok || v == "" {
		return default_
	}
	return v
}

func combineTags(tagParts ...string) []string {
	tags := make([]string, 0)
	for _, element := range tagParts {
		if element != "" {
			tags = append(tags, strings.Split(element, ",")...)
		}
	}
	return tags
}

func serviceMetaData(config *container.Config, port string) (map[string]string, map[string]bool) {

	serviceRegex := regexp.MustCompile("([^_ .]+|^service+)((^[_.]+))?")

	meta := config.Env
	for k, v := range config.Labels {
		meta = append(meta, k+"="+v)
	}
	metadata := make(map[string]string)
	metadataFromPort := make(map[string]bool)
	for _, kv := range meta {
		kvp := strings.SplitN(kv, "=", 2)
		match := serviceRegex.FindAllStringSubmatch(kvp[0],-1)
		
		if len(match)>=1   && strings.EqualFold(match[0][0], "service") {

			key := match[1][0]
			if metadataFromPort[key] {
				continue
			}
			portkey, err := strconv.Atoi(match[1][0])
		//	_, err := strconv.Atoi(portkey[0])
			if err == nil && portkey > 1 {
				if match[1][0] != port {
					continue
				}
				metadata[match[1][0]] = kvp[1]
				metadataFromPort[match[1][0]] = true
			} else {
				metadata[key] = kvp[1]
			}
		}
	}
	return metadata, metadataFromPort
}

func servicePort(container *types.ContainerJSON, port nat.Port, published []nat.PortBinding) DockerServicePort {
	var hp, hip, ep, ept, eip string
	if len(published) > 0 {
		hp = published[0].HostPort
		hip = published[0].HostIP
	}
	if hip == "" {
		hip = "0.0.0.0"
	}
	exposedPort := strings.Split(string(port), "/")
	ep = exposedPort[0]
	if len(exposedPort) == 2 {
		ept = exposedPort[1]
	} else {
		ept = "tcp" // default
	}

	// Nir: support docker NetworkSettings
	eip = container.NetworkSettings.IPAddress
	if eip == "" {
		for _, network := range container.NetworkSettings.Networks {
			eip = network.IPAddress
		}
	}

	return DockerServicePort{
		
		ServicePort: api.ServicePort {
			HostPort:          hp,
			HostIP:            hip,
			ExposedPort:       ep,
			ExposedIP:         eip,
			PortType:          ept,
		},
		ContainerID:       container.ID,
		ContainerHostname: container.Config.Hostname,
		container:         container,
	}
}
