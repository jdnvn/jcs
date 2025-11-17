package main

import (
	"errors"
	"fmt"
	"crypto/rand"
	"encoding/hex"
)

type ServiceHandler struct {
	Services map[string]Service
}

func NewServiceHandler() *ServiceHandler {
    return &ServiceHandler{
        Services: make(map[string]Service),
    }
}

func (s *ServiceHandler) GetService(ID string) (Service, error) {
	service, ok := s.Services[ID]
	if !ok {
		return service, errors.New(fmt.Sprintf("Service not found with ID '%s'", ID))
	}
	return service, nil
}

func (s *ServiceHandler) ListServices() ([]Service, error) {
	services := make([]Service, 0, len(s.Services))
	for _, service := range s.Services {
		services = append(services, service)
	}

	return services, nil
}

func (s *ServiceHandler) CreateService(name string) (Service, error) {
	var newService Service
	for _, service := range s.Services {
		if service.Name == name {
			return newService, errors.New(fmt.Sprintf("A service already exists with the name '%s'", name))
		}
	}

	serviceID, err := s.generateId()
	if err != nil {
		return newService, err
	}
	newService = Service{ID: serviceID, Name: name, Containers: make(map[string]Container)}
	s.Services[serviceID] = newService

	return newService, nil
}

func (s *ServiceHandler) DeleteService(ID string) error {
	service, ok := s.Services[ID]
	if !ok {
		return errors.New(fmt.Sprintf("Service not found with ID '%s'", ID))
	}

	delete(s.Services, service.ID)

	return nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
	  return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *ServiceHandler) generateId() (string, error) {
	for {
		id, err := randomHex(3)
		if err != nil {
			return "", err
		}
		_, ok := s.Services[id]
		if !ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("Could not generate service ID")
}
