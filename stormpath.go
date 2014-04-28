package stormpath

import ()

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

// A list of Stormpath Applications.
type ApplicationList struct {
	Href   string         `json:"href"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Items  *[]Application `json:"items"`
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

// A list of Stormpath Directories.
type DirectoryList struct {
	Href   string       `json:"href"`
	Offset int          `json:"offset"`
	Limit  int          `json:"limit"`
	Items  *[]Directory `json:"items"`
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
