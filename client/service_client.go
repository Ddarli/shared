package client

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type ServiceClient struct {
	consul *consulapi.Client
}

func NewServiceClient(consulAddr string) (*ServiceClient, error) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{consul: client}, nil
}

func (sc *ServiceClient) GetServiceAddress(serviceName string) (string, error) {
	services, _, err := sc.consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("service '%s' not found", serviceName)
	}

	service := services[0].Service
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}
