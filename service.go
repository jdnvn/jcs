package main

import (
	"errors"
	"fmt"
	"net/http"
	"io"
	"log"
	"bytes"
	"encoding/json"
)

type Service struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Containers map[string]Container `json:"containers"`
}

func (s *Service) CreateContainer(imageName string, startCommand string) (Container, error) {
	var newContainer Container

	result, err := serverHandler.ListServers()
	if err != nil {
		return newContainer, err
	}

	var server Server
	if len(result) < 1 {
		randomString, err := randomHex(3)
		if err != nil {
			return newContainer, err
		}
		server, err = serverHandler.CreateServer(fmt.Sprintf("jcs-%s", randomString))
		if err != nil {
			return newContainer, err
		}
	} else {
		server = result[0] // TODO: eventually select a server that has capacity
	}

	sandboxCreateRequest := SandboxCreateRequest{ImageName: imageName, StartCommand: startCommand}

	sandboxCreateRequestJson, err := json.Marshal(sandboxCreateRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return newContainer, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api/sandboxes", server.IP), bytes.NewBuffer(sandboxCreateRequestJson))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return newContainer, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return newContainer, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return newContainer, err
	}

	var sandbox Sandbox
	err = json.Unmarshal(body, &sandbox)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
		return newContainer, err
	}

	containerID, err := randomHex(3)
	if err != nil {
		return newContainer, err
	}

	newContainer = Container{ID: containerID, ServiceID: s.ID, ServerID: server.ID, SandboxID: sandbox.ID, Host: sandbox.PreviewURL, Status: sandbox.Status, ImageName: imageName, StartCommand: startCommand}
	s.Containers[containerID] = newContainer

	return newContainer, nil
}

func (s *Service) ListContainers() ([]Container, error) {
	containers := make([]Container, 0, len(s.Containers))
	for _, container := range s.Containers {
		server, err := serverHandler.GetServer(container.ServerID)
		if err != nil {
			return containers, err
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/sandboxes/%s", server.IP, container.SandboxID), nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
			return containers, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return containers, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
			return containers, err
		}

		var sandbox Sandbox
		err = json.Unmarshal(body, &sandbox)
		if err != nil {
			log.Fatalf("error unmarshalling JSON: %v", err)
			return containers, err
		}
		container.Status = sandbox.Status
	
		containers = append(containers, container)
	}

	return containers, nil
}

func (s *Service) GetContainer(ID string) (Container, error) {
	container, ok := s.Containers[ID]
	if ok {
		return container, errors.New(fmt.Sprintf("Container not found with ID: %s", ID))
	}
	server, err := serverHandler.GetServer(container.ServerID)
	if err != nil {
		return container, err
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/sandboxes/%s", server.IP, container.SandboxID), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return container, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return container, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return container, err
	}

	var sandbox Sandbox
	err = json.Unmarshal(body, &sandbox)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
		return container, err
	}
	container.Status = sandbox.Status

	return container, nil
}
