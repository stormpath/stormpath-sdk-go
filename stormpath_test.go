package stormpath

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     os.Getenv("STORMPATH_API_KEY_ID"),
		Secret: os.Getenv("STORMPATH_API_KEY_SECRET"),
	})
	if err != nil {
		t.Error(err)
	}

	if client.Tenant.Href == "" {
		t.Error("No tenant href could be found.")
	}
}

func TestClientGetTenant(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     os.Getenv("STORMPATH_API_KEY_ID"),
		Secret: os.Getenv("STORMPATH_API_KEY_SECRET"),
	})
	if err != nil {
		t.Error(err)
	}

	tenant, err := client.GetTenant()
	if err != nil {
		t.Error(err)
	}

	if tenant.Href == "" {
		t.Error("No tenant href could be found.")
	} else if tenant.Name == "" {
		t.Error("No tenant name could be found.")
	} else if tenant.Key == "" {
		t.Error("No tenant key could be found.")
	}
}

func TestClientRequest(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     os.Getenv("STORMPATH_API_KEY_ID"),
		Secret: os.Getenv("STORMPATH_API_KEY_SECRET"),
	})
	if err != nil {
		t.Error(err)
	}

	resp, err := client.Request("GET", client.Tenant.Href+"/applications", nil)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Error(err)
	}
}
