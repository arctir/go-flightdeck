package client

import (
	"net/http"

	apiv1 "github.com/arctir/go-flightdeck/pkg/api/v1"
)

func OidcClient(config Config) (*http.Client, error) {
	oauth2Config, err := NewOauthConfig(config.AuthEndpoint)
	if err != nil {
		return nil, err
	}
	oidcConfig := OIDCConfig{
		Config:       *oauth2Config,
		AccessToken:  config.AccessToken,
		RefreshToken: config.RefreshToken,
	}

	return oidcConfig.NewClient()
}

// NewClient provides an authenticated HTTP client for use with the arctir API
func NewClient(endpoint string, config Config) (*apiv1.ClientWithResponses, error) {
	oidcClient, err := OidcClient(config)
	if err != nil {
		return nil, err
	}

	return apiv1.NewClientWithResponses(endpoint, apiv1.WithHTTPClient(oidcClient))
}
