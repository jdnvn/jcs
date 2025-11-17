package main

import (
	"log"
	"fmt"
	"github.com/joho/godotenv"
	"errors"
)

type ServerHandler struct {
	Servers map[string]Server
	ServerAdapter ServerAdapter
}

func NewServerHandler() *ServerHandler {
	servers := make(map[string]Server)
	// instead of this, we should get the API key from a "config"
	err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file")
    }
	serverAdapter := LocalServerAdapter{} // TODO: set this based on the user's config settings
	remoteServers, err := serverAdapter.ListServers()
	if err == nil {
		for _, remoteServer := range remoteServers {
			id, _ := randomHex(3)
			serverID := fmt.Sprintf("%s%s", remoteServer.ID, id)
			server := Server{ID: serverID, RemoteID: remoteServer.ID, Name: remoteServer.Name, Type: remoteServer.Type, Status: remoteServer.Status, IP: remoteServer.IP}
			servers[serverID] = server
		}
	}

    return &ServerHandler{
        Servers: servers,
		ServerAdapter: serverAdapter,
    }
}

func (s *ServerHandler) GetServer(ID string) (Server, error) {
	server, ok := s.Servers[ID]
	if !ok {
		return server, errors.New(fmt.Sprintf("Server not found with ID '%s'", ID))
	}
	return server, nil
}

func (s *ServerHandler) ListServers() ([]Server, error) {
	servers := make([]Server, 0, len(s.Servers))
	for _, server := range s.Servers {
		servers = append(servers, server)
	}

	return servers, nil
}

func (s *ServerHandler) CreateServer(name string) (Server, error) {
	var newServer Server
	for _, server := range s.Servers {
		if server.Name == name {
			return newServer, errors.New(fmt.Sprintf("A server already exists with the name '%s'", name))
		}
	}

	remoteServer, err := s.ServerAdapter.CreateServer(name)
	if err != nil {
		return newServer, err
	}

	id, err := randomHex(3)
	if err != nil {
		return newServer, err
	}
	serverID := fmt.Sprintf("%s%s", remoteServer.ID, id)
	newServer = Server{ID: serverID, RemoteID: remoteServer.ID, Name: remoteServer.Name, Type: remoteServer.Type, Status: remoteServer.Status, IP: remoteServer.IP}
	s.Servers[newServer.ID] = newServer

	return newServer, nil
}

func (s *ServerHandler) DeleteServer(ID string) error {
	server, ok := s.Servers[ID]
	if !ok {
		return errors.New(fmt.Sprintf("Server not found with ID '%s'", ID))
	}

	// TODO: delete from remote?
	delete(s.Servers, server.ID)

	return nil
}

func (s *ServerHandler) generateId() (string, error) {
	for {
		id, err := randomHex(3)
		if err != nil {
			return "", err
		}
		_, ok := s.Servers[id]
		if !ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("Could not generate server ID")
}
