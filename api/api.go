package api

import (
)

type EventProcessor func(status string, service *Service, closeChan chan error) error

type RegistryAdapter interface {
    Ping() error
    RunTemplate(status string, service *Service) error

}

type Service struct {
	ID    string
	Name  string
	Port  int
	IP    string
	Version string
	Tags  []string
	Attrs map[string]string
	TTL   int
	Container struct{}

	Origin ServicePort
}

type ServicePortAPI interface {
	getHostPort() string
	getHostIp() string
	getExposedPort() string
	getExposedIp() string
	getPortType() string
}


type ServicePort struct {
	HostPort          string
	HostIP            string
	ExposedPort       string
	ExposedIP         string
	PortType          string
}

func (s *ServicePort) getHostPort() string {
	return s.HostPort
}

func (s *ServicePort) getHostIp() string {
	return s.HostIP
}

func (s *ServicePort) getExposedPort() string {
	return s.ExposedPort
}

func (s *ServicePort) getExposedIp() string {
	return s.ExposedIP
}

func (s *ServicePort) getPortType() string {
	return s.PortType
}
