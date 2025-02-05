package coze

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// User represents a Coze user
type User struct {
	baseModel
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	AvatarURL string `json:"avatar_url"`
}

// UsersClient represents the client for users-related operations
type UsersClient struct {
	core     *core
	baseURL  string
	hostName string
}

type usersOption struct {
	baseURL    string
	httpClient HTTPClient
}

type UsersClientOption func(*usersOption)

// WithUsersBaseURL adds base URL
func WithUsersBaseURL(baseURL string) UsersClientOption {
	return func(opt *usersOption) {
		opt.baseURL = baseURL
	}
}

// WithUsersHttpClient adds HTTP client
func WithUsersHttpClient(client HTTPClient) UsersClientOption {
	return func(opt *usersOption) {
		opt.httpClient = client
	}
}

// NewUsersClient creates a new users client
func NewUsersClient(opts ...UsersClientOption) (*UsersClient, error) {
	initSettings := &usersOption{
		baseURL: ComBaseURL,
	}

	for _, opt := range opts {
		opt(initSettings)
	}

	var hostName string
	if initSettings.baseURL != "" {
		parsedURL, err := url.Parse(initSettings.baseURL)
		if err != nil {
			return nil, err
		}
		hostName = parsedURL.Host
	} else {
		return nil, errors.New("base URL is required")
	}

	var httpClient HTTPClient
	if initSettings.httpClient != nil {
		httpClient = initSettings.httpClient
	} else {
		httpClient = http.DefaultClient
	}

	return &UsersClient{
		core:     newCore(httpClient, initSettings.baseURL),
		baseURL:  initSettings.baseURL,
		hostName: hostName,
	}, nil
}

// Me retrieves the current user's information
func (c *UsersClient) Me(ctx context.Context) (*User, error) {
	result := &User{}
	if err := c.core.Request(ctx, http.MethodGet, "/v1/users/me", nil, result); err != nil {
		return nil, err
	}
	return result, nil
}
