package client

import (
	"context"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const defaultClientID = "arctir-cli"

func NewOauthConfig(endpoint string) (*oauth2.Config, error) {
	scopes := []string{"profile", "email", "openid", "group", "offline_access"}
	provider, err := oidc.NewProvider(context.Background(), endpoint)
	if err != nil {
		return nil, err
	}
	return &oauth2.Config{
		ClientID: defaultClientID,
		Endpoint: provider.Endpoint(),
		Scopes:   scopes,
	}, nil
}
