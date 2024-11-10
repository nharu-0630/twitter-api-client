package api

import "github.com/nharu-0630/twitter-api-client/model"

type ClientsPipe struct {
	clients []Client
}

func NewClientsPipe(config model.Config) *ClientsPipe {
	var clients []Client
	for _, token := range config.Tokens {
		clients = append(clients, *NewAuthorizedClient(token.AuthToken, token.CsrfToken))
	}
	return &ClientsPipe{
		clients: clients,
	}
}

func (cp *ClientsPipe) operation(operation string) Client {
	for _, client := range cp.clients {
		if _, ok := client.rateLimits[operation]; ok {
			if client.rateLimits[operation].remaining == 0 {
				continue
			}
			return client
		}
	}
	var minReset int
	var minClient Client
	for _, client := range cp.clients {
		if _, ok := client.rateLimits[operation]; ok {
			if minReset == 0 || client.rateLimits[operation].reset < minReset {
				minReset = client.rateLimits[operation].reset
				minClient = client
			}
		}
	}
	return minClient
}
