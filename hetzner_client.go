package main

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"log"
	"bytes"
	"encoding/json"
)

type HetznerIPV4Response struct {
	IP string `json:"ip"`
}

type HetznerPublicNetResponse struct {
	IPV4 HetznerIPV4Response `json:"ipv4"`
}

type HetznerServer struct {
	ID int `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	PublicNet HetznerPublicNetResponse `json:"public_net"`
}

type HetznerListServersResponse struct {
	Servers []HetznerServer `json:"servers"`
}

type HetznerCreateServerResponse struct {
	Server HetznerServer `json:"server"`
}

type HetznerGetServerResponse struct {
	Server HetznerServer `json:"server"`
}

type HetznerCreateServerRequest struct {
	Name string `json:"name"`
	ServerType string `json:"server_type"`
	Image string `json:"image"`
}

type HetznerApiClient struct {
}

func (api *HetznerApiClient) GetServer(serverID string) (HetznerGetServerResponse, error) {
	var result HetznerGetServerResponse
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.hetzner.cloud/v1/servers/%s", serverID), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("HETZNER_API_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
		return result, err
	}

	return result, nil
}

func (api *HetznerApiClient) ListServers() (HetznerListServersResponse, error) {
	var result HetznerListServersResponse

    apiKey := os.Getenv("HETZNER_API_KEY")
    if apiKey == "" {
        log.Fatalf("HETZNER_API_KEY not set!")
        return result, fmt.Errorf("HETZNER_API_KEY not set")
    }

	req, err := http.NewRequest(http.MethodGet, "https://api.hetzner.cloud/v1/servers", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return result, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
		return result, err
	}

	return result, nil
} 

func (api *HetznerApiClient) CreateServer(name string, serverType string, image string) (HetznerCreateServerResponse, error) {
	var result HetznerCreateServerResponse
	requestBody := HetznerCreateServerRequest{Name: name, ServerType: serverType, Image: image}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return result, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.hetzner.cloud/v1/servers", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return result, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("HETZNER_API_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
		return result, err
	}

	return result, nil
}
