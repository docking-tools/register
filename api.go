package main

import (
    dockerapi "github.com/fsouza/go-dockerclient"
)

type RegistryAdapter interface {
    Ping() error
    Register(service *Service) error
    Deregsiter(service *Service) error
    Update(service *Service) error
}

type Service struct {
	ID    string
	Name  string
	Port  int
	IP    string
	Tags  []string
	Attrs map[string]string
	TTL   int

	Origin ServicePort
}

type ServicePort struct {
	HostPort          string
	HostIP            string
	ExposedPort       string
	ExposedIP         string
	PortType          string
	ContainerHostname string
	ContainerID       string
	ContainerName     string
	container         *dockerapi.Container
}