package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-connections/nat"
	api "github.com/docking-tools/register/api"
	"regexp"
	"strconv"
	"strings"
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

	serviceRegex := regexp.MustCompile("([^_.]+|^service[_.]+)((^[_.]+))?")

	meta := config.Env
	for k, v := range config.Labels {
		meta = append(meta, k+"="+v)
	}
	metadata := make(map[string]string)
	metadataFromPort := make(map[string]bool)
	for _, kv := range meta {
		kvp := strings.SplitN(kv, "=", 2)
		match := serviceRegex.FindAllStringSubmatch(kvp[0], -1)

		if len(match) >= 1 && strings.EqualFold("service", match[0][0]) {

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
				keys := make([]string, 0)
				for toto := range match[2:] {

					keys = append(keys, match[2+toto][0])
				}
				metadata[strings.ToLower(strings.Join(keys, "."))] = kvp[1]
				metadataFromPort[strings.ToLower(strings.Join(keys, "."))] = true

			} else {
				keys := make([]string, 0)
				for toto := range match[1:] {

					keys = append(keys, match[1+toto][0])
				}
				metadata[strings.ToLower(strings.Join(keys, "."))] = kvp[1]
			}
		}
	}
	return metadata, metadataFromPort
}

func graphMetaData(config *container.Config) api.Recmap {
	meta := config.Env
	for k, v := range config.Labels {
		meta = append(meta, k+"="+v)
	}
	metaRegex := regexp.MustCompile("[_.]")

	//var nextMap interface{}
	nextMap := make(api.Recmap)
	result := nextMap
	for _, kv := range meta {
		kvp := strings.SplitN(kv, "=", 2)
		if len(kvp) >= 2 {
			match := metaRegex.Split(kvp[0], -1)
			for _, key := range match[:len(match)-1] {
				sKey := strings.ToLower(key)
				if _, ok := nextMap[sKey].(api.Recmap); !ok {
					nextMap[sKey] = make(api.Recmap)
				}
				nextMap = nextMap[sKey].(api.Recmap)

			}
			nextMap[match[len(match)-1]] = kvp[1]
			nextMap = result
		}
	}
	return result

}

func serviceSwarmPort(container *types.ContainerJSON, port swarm.PortConfig) DockerServicePort {
	var hp, ep, ept string

	hp = string(port.TargetPort)
	ep = string(port.PublishedPort)
	ept = string(port.PublishMode)

	return DockerServicePort{

		ServicePort: api.ServicePort{
			HostPort:    hp,
			HostIP:      "",
			ExposedPort: ep,
			ExposedIP:   "",
			PortType:    ept,
		},
		ContainerID:       container.ID,
		ContainerHostname: container.Config.Hostname,
		container:         container,
	}
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

		ServicePort: api.ServicePort{
			HostPort:    hp,
			HostIP:      hip,
			ExposedPort: ep,
			ExposedIP:   eip,
			PortType:    ept,
		},
		ContainerID:       container.ID,
		ContainerHostname: container.Config.Hostname,
		container:         container,
	}
}
