package stormpath

import (
	"github.com/nu7hatch/gouuid"
	"log"
	"os"
	"testing"
)

var (
	STORMPATH_API_KEY_ID     = os.Getenv("STORMPATH_API_KEY_ID")
	STORMPATH_API_KEY_SECRET = os.Getenv("STORMPATH_API_KEY_SECRET")
	TEST_PREFIX              string
	CLIENT                   *Client
)

func init() {
	if STORMPATH_API_KEY_ID == "" {
		log.Fatal("STORMPATH_API_KEY_ID not set in the environment.")
	} else if STORMPATH_API_KEY_SECRET == "" {
		log.Fatal("STORMPATH_API_KEY_SECRET not set in the environment.")
	}

	// Generate a globally unique UUID to be used as a prefix throughout our
	// testing.
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("UUID generation failed.")
	}

	// Store our test prefix.
	TEST_PREFIX = uuid.String()

	// Generate a Stormpath client we'll use for all our tests.
	client, err := NewClient(&ApiKeyPair{
		Id:     STORMPATH_API_KEY_ID,
		Secret: STORMPATH_API_KEY_SECRET,
	})
	if err != nil {
		log.Fatal("Couldn't create a Stormpath client.")
	}
	CLIENT = client
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

func TestRequest(t *testing.T) {
	resp, err := CLIENT.Request("GET", "/tenants/current", nil)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		t.Error(err)
	}
}

func TestGetTenant(t *testing.T) {
	tenant, err := CLIENT.GetTenant()
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

func TestGetApplications(t *testing.T) {
	_, err := CLIENT.GetApplications()
	if err != nil {
		t.Error(err)
	}
}

func TestGetDirectories(t *testing.T) {
	_, err := CLIENT.GetDirectories()
	if err != nil {
		t.Error(err)
	}
}
