package main

import (
	"errors"
	"fmt"
)

type LocalServerAdapter struct {

}

func (l LocalServerAdapter) ListServers() ([]RemoteServer, error) {
	result := []RemoteServer {
		RemoteServer{ID: "localhost", Name: "localhost", Type: "local", Status: "online", IP: "localhost"},
	}

	return result, nil
}

func (l LocalServerAdapter) GetServer(id string) (RemoteServer, error) {
	var result RemoteServer
	if id != "localhost" {
		return result, errors.New(fmt.Sprintf("Server not found with ID: '%s'", id))
	}
	result = RemoteServer{ID: "localhost", Name: "localhost", Type: "local", Status: "online", IP: "localhost"}

	return result, nil
}

func (l LocalServerAdapter) CreateServer(name string) (RemoteServer, error) {
	result := RemoteServer{ID: fmt.Sprintf("localhost-%s", name), Name: name, Type: "local", Status: "online", IP: "localhost"}

	return result, nil
}
