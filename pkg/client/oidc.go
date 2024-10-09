package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	Config       oauth2.Config
	AccessToken  string
	RefreshToken string
}

func ExtractExpiry(idToken string, deltaSeconds int) (*time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not extract token claims")
	}

	var tm time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		tm = time.Unix(v, 0)
	}

	tm = tm.Add(time.Duration((-1 * deltaSeconds)) * time.Second)
	return &tm, nil
}

func (c OIDCConfig) NewClient() (*http.Client, error) {

	exp, err := ExtractExpiry(c.AccessToken, 10)
	if err != nil {
		return nil, err
	}

	token := oauth2.Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		Expiry:       *exp,
	}

	ts := c.Config.TokenSource(context.Background(), &token)
	rts := oauth2.ReuseTokenSource(&token, ts)

	baseTransport := oidcHTTPTransport{
		T:   http.DefaultTransport,
		rts: rts,
	}

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: rts,
			Base:   &baseTransport,
		},
	}
	return client, nil
}

type oidcHTTPTransport struct {
	T   http.RoundTripper
	rts oauth2.TokenSource
}

func (t *oidcHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, _ := t.rts.Token()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))
	return t.T.RoundTrip(req)
}
