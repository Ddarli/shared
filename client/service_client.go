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

func (sc *ServiceClient) RegisterService(serviceName string, address string, port int) error {
	// Создаем регистрацию сервиса
	registration := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, address, port),
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
			Interval: "10s",
			Timeout:  "5s",
		},
		Tags: []string{"microservice", "api"},
	}

	// Регистрируем сервис
	if err := sc.consul.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	return nil
}

func (sc *ServiceClient) DeregisterService(serviceID string) error {
	if err := sc.consul.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	return nil
}
