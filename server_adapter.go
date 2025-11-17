package main

type ServerAdapter interface {
	ListServers() ([]RemoteServer, error)
	GetServer(id string) (RemoteServer, error)
	CreateServer(name string) (RemoteServer, error)
}
