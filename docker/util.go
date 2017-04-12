package docker

import (
	"github.com/docker/docker/api/types"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-connections/nat"
	api "github.com/docking-tools/register/api"
	"path"
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

func envArrayToMap(envs []string) map[string]string {
	result := make(map[string]string)
	for _, kv := range envs {
		kvp := strings.SplitN(kv, "=", 2)
		if len(kvp) >= 2 {
			result[kvp[0]] = kvp[1]
		}
	}
	return result
}

func serviceMetaData(port string, metadata ...map[string]string) (map[string]string, map[string]bool) {

	serviceRegex := regexp.MustCompile("([^_.]+|^service[_.]+)((^[_.]+))?")

	meta := make(map[string]string)
	for _, arg := range metadata {
		for k, v := range arg {
			meta[k] = v
		}
	}
	serviceMetaData := make(map[string]string)
	metadataFromPort := make(map[string]bool)
	for k, v := range meta {
		match := serviceRegex.FindAllStringSubmatch(k, -1)

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
				serviceMetaData[strings.ToLower(strings.Join(keys, "."))] = v
				metadataFromPort[strings.ToLower(strings.Join(keys, "."))] = true

			} else {
				keys := make([]string, 0)
				for toto := range match[1:] {

					keys = append(keys, match[1+toto][0])
				}
				serviceMetaData[strings.ToLower(strings.Join(keys, "."))] = v
			}
		}
	}
	return serviceMetaData, metadataFromPort
}

func graphMetaData(metadata ...map[string]string) api.Recmap {
	meta := make(map[string]string)
	for _, arg := range metadata {
		for k, v := range arg {
			meta[k] = v
		}
	}
	metaRegex := regexp.MustCompile("[_.]")

	//var nextMap interface{}
	nextMap := make(api.Recmap)
	result := nextMap
	for k, v := range meta {
		match := metaRegex.Split(k, -1)
		for _, key := range match[:len(match)-1] {
			sKey := strings.ToLower(key)
			if _, ok := nextMap[sKey].(api.Recmap); !ok {
				nextMap[sKey] = make(api.Recmap)
			}
			nextMap = nextMap[sKey].(api.Recmap)

		}
		nextMap[match[len(match)-1]] = v
		nextMap = result
	}
	return result

}

func serviceSwarmPort(container *types.ContainerJSON, swarmService *swarm.Service, port swarm.PortConfig) DockerServicePort {
	var hp, ep, ept string

	ep = strconv.FormatUint(uint64(port.TargetPort), 10)
	hp = strconv.FormatUint(uint64(port.PublishedPort), 10)
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
		ServiceName:       swarmService.Spec.Name,
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
		ServiceName:       strings.Split(path.Base(container.Config.Image), ":")[0],
	}
}
