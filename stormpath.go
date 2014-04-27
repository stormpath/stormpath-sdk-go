package stormpath

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	VERSION    = "0.0.1"
	USER_AGENT = "go-stormpath/" + VERSION
	BASE_URL   = "https://api.stormpath.com/v1"
	ENABLED    = "enabled"
	DISABLED   = "disabled"
)

// The ApiKeyPair object is meant for storing Stormpath credentials.
type ApiKeyPair struct {
	Id     string
	Secret string
}

// The Client object allows you to easily communicate with the Stormpath API
// service.
type Client struct {
	Keypair   *ApiKeyPair
	Transport *http.Transport
	Tenant    *Tenant
}

// A Tenant is a globally unique namespace.
type Tenant struct {
	Href string `json:"href"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// An Application is a unique Stormpath application.
type Application struct {
	Href        string  `json:"href"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Tenant      *Tenant `json:"-"`
}

// A generic Stormpath collection -- this is used internally to handle
// deserializing large collections of resources.
type Collection struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// An ApplicationCollection is a list of Applications.
type ApplicationCollection struct {
	Href   string         `json:"href"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Items  *[]Application `json:"items"`
}

// A unique Stormpath directory.
type Directory struct {
	Href        string
	Name        string
	Description string
	Status      string
	Tenant      *Tenant
}

// A unique Stormpath group.
type Group struct {
	Href        string
	Name        string
	Description string
	Status      string
	Tenant      *Tenant
	Directory   *Directory
}

// A unique Stormpath group membership.
type GroupMembership struct {
	Href    string
	Account *Account
	Group   *Group
}

// A unique Stormpath account.
type Account struct {
	Href       string
	Username   string
	Email      string
	Password   string
	FullName   string
	GivenName  string
	MiddleName string
	Surname    string
	Status     string
	Groups     *GroupMembership
	Directory  *Directory
	Tenant     *Tenant
}

// Custom data.
type CustomData struct {
	Href       string
	CreatedAt  string
	ModifiedAt string
}

// NewClient generates a new Client, given an ApiKeyPair as input.  It will then
// use the ApiKeyPair to attempt to fetch the current Stormpath Tenant.  Returns
// an initialized Client (thread safe) and any error encountered.
func NewClient(keypair *ApiKeyPair) (*Client, error) {
	client := &Client{
		Transport: &http.Transport{},
		Keypair:   keypair,
	}

	tenant, err := client.GetTenant()
	if err != nil {
		return nil, err
	}
	client.Tenant = tenant

	return client, nil
}

// Request makes an HTTP request to the Stormpath API, using the user's
// credentials automatically.  Returns an HTTP response and any error
// encountered.
func (client *Client) Request(method string, url string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("User-Agent", USER_AGENT)
	req.SetBasicAuth(client.Keypair.Id, client.Keypair.Secret)

	resp, err := client.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTenant returns the current Tenant and any error encountered.
func (client *Client) GetTenant() (*Tenant, error) {
	resp, err := client.Request("GET", BASE_URL+"/tenants/current", nil)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode != 302 {
		return nil, errors.New("Didn't receive HTTP 302 when requesting tenant.")
	} else if resp.Header.Get("Location") == "" {
		return nil, errors.New("No Location header found when requesting tenant.")
	}

	resp, err = client.Request("GET", resp.Header.Get("Location"), nil)
	if err != nil {
		return nil, errors.New("Got HTTP error when requesting tenant.")
	}
	defer resp.Body.Close()

	tenant := &Tenant{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

// GetApplications returns an ApplicationCollection, and any error encountered.
// TODO: Iterate over *all* applications, not just some.
func (client *Client) GetApplications() (*ApplicationCollection, error) {
	resp, err := client.Request("GET", client.Tenant.Href+"/applications", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//collection := &Collection{}
	ac := &ApplicationCollection{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(ac)
	if err != nil {
		return nil, err
	}

	return ac, nil
}
