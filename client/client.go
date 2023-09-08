package client

type APIClient struct {
	Active             bool
	ClientID           string
	ApiKey             string
	SoftwareProviderID int
	Sandbox            bool
	Description        string
}
