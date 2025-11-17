package main

import (
	"strconv"
)

type HetznerServerAdapter struct {
}

func (h HetznerServerAdapter) ListServers() ([]RemoteServer, error) {
	result := []RemoteServer{}
	api := HetznerApiClient{}
	listServersResponse, err := api.ListServers()
	if err != nil {
		return result, err
	}
	for _, serverData := range listServersResponse.Servers {
		server := RemoteServer{ID: strconv.Itoa(serverData.ID), Type: "hetzner", Status: serverData.Status, IP: serverData.PublicNet.IPV4.IP}
		result = append(result, server)
	}
	return result, nil
}

func (h HetznerServerAdapter) GetServer(ID string) (RemoteServer, error) {
	var result RemoteServer
	api := HetznerApiClient{}
	getServerResponse, err := api.GetServer(ID)
	if err != nil {
		return result, err
	}
	hetznerServer := getServerResponse.Server
	result = RemoteServer{ID: strconv.Itoa(hetznerServer.ID), Name: hetznerServer.Name, Type: "hetzner", Status: hetznerServer.Status, IP: hetznerServer.PublicNet.IPV4.IP}
	return result, nil
}

func (h HetznerServerAdapter) CreateServer(name string) (RemoteServer, error) {
	var result RemoteServer
	api := HetznerApiClient{}
	createServerResponse, err := api.CreateServer(name, "cpx21", "ubuntu-24.04")
	if err != nil {
		return result, err
	}
	hetznerServer := createServerResponse.Server
	result = RemoteServer{ID: strconv.Itoa(hetznerServer.ID), Name: hetznerServer.Name, Type: "hetzner", Status: hetznerServer.Status, IP: hetznerServer.PublicNet.IPV4.IP}
	return result, nil
}
