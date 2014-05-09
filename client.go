package stormpath

import (
	"bytes"
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
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

// CreateApplication creates a new Stormpath Application resource for the given
// Tenant.  Returns a new Application object, and any error encountered.
func (client *Client) CreateApplication(application *Application, createDirectory bool) (*Application, error) {

	// First, convert the Application to JSON.
	jsonBytes, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}

	// Fire off the creation request.
	resp, err := client.Request("POST", "/applications", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	dec := json.NewDecoder(resp.Body)

	// If the response didn't generate a 201 CREATED, this means something bad
	// happened, and we can expect an error from Stormpath.
	if resp.StatusCode != 201 {
		se := &StormpathError{}

		err = dec.Decode(se)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("%s More Info: %s", se.DeveloperMessage, se.MoreInfo))
	}

	// If we get here, it means the Application was successfully created! So
	// we'll then create a new Application object, and return it.
	newApp := &Application{}

	err = dec.Decode(newApp)
	if err != nil {
		return nil, err
	}

	return newApp, nil
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

// GetDirectories gets a list of all Stormpath Directories for the given
// Tenant.  Returns a slice of Directories and any error encountered.
// NOTE: This may take a while, as if you have a lot of Directories, this will
// iterate over them all before returning.
func (client *Client) GetDirectories() (*[]Directory, error) {
	resp, err := client.Request("GET", client.Tenant.Href+"/directories", nil)
	if err != nil {
		return nil, err
	}

	dl := &DirectoryList{}
	err = json.NewDecoder(resp.Body).Decode(&dl)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// If there is only a single page of directories, return immediately.
	if len(*dl.Items) < dl.Limit {
		return dl.Items, nil
	}

	// Create a new slice called dirs -- this will hold our final list of
	// Directory objects that are eventually returned.
	dirs := make([]Directory, len(*dl.Items))
	copy(dirs, *dl.Items)

	// Loop through all subsequent pages of Directories, adding each
	// Directory to our final dirs slice.
	for offset := dl.Limit; len(*dl.Items) == dl.Limit; offset += dl.Limit {
		resp, err := client.Request("GET", fmt.Sprintf("%v/directories?offset=%v", client.Tenant.Href, offset), nil)
		if err != nil {
			return nil, err
		}

		// Grab the JSON body of the request.
		dl = &DirectoryList{}
		err = json.NewDecoder(resp.Body).Decode(&dl)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()

		// Append each Directory onto our dirs slice.
		for _, dir := range *dl.Items {
			dirs = append(dirs, dir)
		}
	}

	return &dirs, nil
}
