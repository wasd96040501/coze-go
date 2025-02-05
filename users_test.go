package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersClient_Me(t *testing.T) {
	mockTransport := &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			expectedUser := &meResp{
				User: &User{
					UserID:    "test_user_id",
					UserName:  "test_user",
					NickName:  "Test User",
					AvatarURL: "https://example.com/avatar.jpg",
				},
			}
			return mockResponse(http.StatusOK, expectedUser)
		},
	}

	client := NewCozeAPI(NewTokenAuth("test_token"),
		WithBaseURL(ComBaseURL),
		WithHttpClient(&http.Client{Transport: mockTransport}),
	)

	user, err := client.Users.Me(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "test_user_id", user.UserID)
	assert.Equal(t, "test_user", user.UserName)
	assert.Equal(t, "Test User", user.NickName)
	assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
}
