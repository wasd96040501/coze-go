package coze

import (
	"context"
	"net/http"
)

// User represents a Coze user
type User struct {
	baseModel
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	AvatarURL string `json:"avatar_url"`
}

type users struct {
	client *core
}

func newUsers(core *core) *users {
	return &users{
		client: core,
	}
}

// Me retrieves the current user's information
func (r *users) Me(ctx context.Context) (*User, error) {
	method := http.MethodGet
	uri := "/v1/users/me"
	result := &User{}
	if err := r.client.Request(ctx, method, uri, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}
