package stormpath

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// The Client object allows you to easily communicate with the Stormpath API
// service.
type Client struct {
	Keypair   *ApiKeyPair
	Transport *http.Transport
	Tenant    *Tenant
}

// NewClient generates a new Stormpath Client, given an ApiKeyPair as input.
// It will then use the ApiKeyPair to attempt to fetch the current Stormpath
// Tenant.  Returns an initialized Client (thread safe) and any error
// encountered.
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

	// If the URL starts with /, we'll go ahead and generate the fully qualified
	// URL string to make handling requests easier.
	if strings.HasPrefix(url, "/") {
		url = BASE_URL + url
	}

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

// GetApplications gets a list of all Stormpath Applications for the given
// Tenant.  Returns a slice of Applications and any error encountered.
// NOTE: This may take a while, as if you have a lot of Applications, this will
// iterate over them all before returning.
func (client *Client) GetApplications() (*[]Application, error) {
	resp, err := client.Request("GET", client.Tenant.Href+"/applications", nil)
	if err != nil {
		return nil, err
	}

	al := &ApplicationList{}
	err = json.NewDecoder(resp.Body).Decode(&al)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// If there is only a single page of applications, return immediately.
	if len(*al.Items) < al.Limit {
		return al.Items, nil
	}

	// Create a new slice called apps -- this will hold our final list of
	// Application objects that are eventually returned.
	apps := make([]Application, len(*al.Items))
	copy(apps, *al.Items)

	// Loop through all subsequent pages of Applications, adding each
	// Application to our final apps slice.
	for offset := al.Limit; len(*al.Items) == al.Limit; offset += al.Limit {
		resp, err := client.Request("GET", fmt.Sprintf("%v/applications?offset=%v", client.Tenant.Href, offset), nil)
		if err != nil {
			return nil, err
		}

		// Grab the JSON body of the request.
		al = &ApplicationList{}
		err = json.NewDecoder(resp.Body).Decode(&al)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()

		// Append each Application onto our apps slice.
		for _, app := range *al.Items {
			apps = append(apps, app)
		}
	}

	return &apps, nil
}
