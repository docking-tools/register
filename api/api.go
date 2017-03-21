package api

import (
)

type EventProcessor func(status string, object interface{}, closeChan chan error) error

type Recmap map[string]interface{}
type RegistryAdapter interface {
    RunTemplate(status string, object interface{}) error

}


type Instance struct {
	Container interface{}
	Services []*Service
	MetaDataGraph map[string]interface{}
	Status string
	
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
