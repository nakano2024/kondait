package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"kondait-backend/application/auth"
	"kondait-backend/infra/config"
	"kondait-backend/web/dto"
)

type authIntrospector struct {
	config     config.Config
	httpClient *http.Client
}

func (introspector *authIntrospector) Introspect(ctx context.Context, token string) (auth.AuthIntrospectionResult, error) {
	if introspector.httpClient == nil {
		return auth.AuthIntrospectionResult{}, errors.New("authIntrospector: http client is nil")
	}

	endpoint := strings.TrimRight(introspector.config.AuthServerUrl, "/") +
		"/protocol/openid-connect/token/introspect"
	form := url.Values{}
	form.Set("token", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return auth.AuthIntrospectionResult{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	credentials := introspector.config.ClientId + ":" + introspector.config.ClientSecret
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	resp, err := introspector.httpClient.Do(req)
	if err != nil {
		return auth.AuthIntrospectionResult{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return auth.AuthIntrospectionResult{}, readErr
		}
		return auth.AuthIntrospectionResult{}, errors.New(string(bodyBytes))
	}

	var introspectionResponse struct {
		Active bool   `json:"active"`
		Sub    string `json:"sub"`
		Scope  string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&introspectionResponse); err != nil {
		return auth.AuthIntrospectionResult{}, err
	}

	return auth.AuthIntrospectionResult{
		IsActive: introspectionResponse.Active,
		Sub:      introspectionResponse.Sub,
		Scopes:   strings.Fields(introspectionResponse.Scope),
	}, nil
}

func NewAuthIntrospector(cfg config.Config, httpClient *http.Client) auth.IAuthIntrospector {
	return &authIntrospector{
		config:     cfg,
		httpClient: httpClient,
	}
}

type authIntrospectorMock struct{}

func (introspector *authIntrospectorMock) Introspect(ctx context.Context, token string) (auth.AuthIntrospectionResult, error) {
	_ = ctx
	return auth.AuthIntrospectionResult{
		IsActive: true,
		Sub:      "fac0fa00-7ee9-b423-813f-eee8e115ca17",
		Scopes: []string{
			dto.ScopeCookingItemsRead,
		},
	}, nil
}

func NewAuthIntrospectorMock() auth.IAuthIntrospector {
	return &authIntrospectorMock{}
}
