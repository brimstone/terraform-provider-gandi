package main

import (
	"log"

	"github.com/prasmussen/gandi-api/client"
)

// Config contains DNSMadeEasy provider settings
type Config struct {
	Key     string
	Testing bool
}

// Env gets appropriate system type
func (c *Config) Env() client.SystemType {
	if c.Testing {
		return client.Testing
	}
	return client.Production
}

// Client returns a new client for accessing Gandi API via meta passed to CRUD
func (c *Config) Client() *client.Client {

	gandiClient := client.New(c.Key, c.Env())
	log.Printf("[INFO] Gandi Client configured for URL: %s with Key: %s", gandiClient.Url, c.Key)

	return gandiClient
}
