package skyeng

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"
)

const (
	getUserIDCookieName     = "token_global"
	getUserIDJWTUserIDClaim = "userId"
)

type GetUserIDSpec struct {
	Cookies []*http.Cookie
}

type GetUserIDData struct {
	Cookies []*http.Cookie
	UserID  int
}

func (r *Repository) GetUserID(_ context.Context, spec GetUserIDSpec) (*GetUserIDData, error) {
	var tokenGlobakCookie *http.Cookie

	for _, cookie := range spec.Cookies {
		if cookie.Name == getUserIDCookieName {
			tokenGlobakCookie = cookie
		}
	}

	if tokenGlobakCookie == nil {
		return nil, errors.New("token_global cookie not found")
	}

	userID, err := httphelper.ParseJWTClaimsInt(tokenGlobakCookie.Value, getUserIDJWTUserIDClaim)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse user_id from token_global cookie")
	}

	return &GetUserIDData{
		UserID:  userID,
		Cookies: spec.Cookies,
	}, nil
}
