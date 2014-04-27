package stormpath

import (
	"log"
	"os"
	"testing"
)

var (
	STORMPATH_API_KEY_ID     = os.Getenv("STORMPATH_API_KEY_ID")
	STORMPATH_API_KEY_SECRET = os.Getenv("STORMPATH_API_KEY_SECRET")
)

func init() {
	if STORMPATH_API_KEY_ID == "" {
		log.Fatal("STORMPATH_API_KEY_ID not set in the environment.")
	} else if STORMPATH_API_KEY_SECRET == "" {
		log.Fatal("STORMPATH_API_KEY_SECRET not set in the environment.")
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     STORMPATH_API_KEY_ID,
		Secret: STORMPATH_API_KEY_SECRET,
	})
	if err != nil {
		t.Error(err)
	}

	if client.Tenant.Href == "" {
		t.Error("No tenant href could be found.")
	}
}

func TestClientRequest(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     STORMPATH_API_KEY_ID,
		Secret: STORMPATH_API_KEY_SECRET,
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

func TestClientGetTenant(t *testing.T) {
	client, err := NewClient(&ApiKeyPair{
		Id:     STORMPATH_API_KEY_ID,
		Secret: STORMPATH_API_KEY_SECRET,
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
