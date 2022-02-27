package java_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type authKey struct{}

var (
	ErrMissingJwtToken        = errors.Unauthorized("UNAUTHORIZED", "Oauth2 token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized("UNAUTHORIZED", "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized("UNAUTHORIZED", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("UNAUTHORIZED", "Oauth2 token has expired")
	ErrTokenParseFail         = errors.Unauthorized("UNAUTHORIZED", "Fail to parse Oauth2 token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("UNAUTHORIZED", "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("UNAUTHORIZED", "Wrong context for middleware")
	ErrNeedTokenProvider      = errors.Unauthorized("UNAUTHORIZED", "Token provider is missing")
	ErrSignToken              = errors.Unauthorized("UNAUTHORIZED", "Can not sign token.Is the key correct?")
	ErrGetKey                 = errors.Unauthorized("UNAUTHORIZED", "Can not get key while signing token")
)

const (
	AuthorizationKey = "Authorization"
	BearerWord       = "Bearer"
)

type Response struct {
	License  string `json:"license"`
	UserInfo struct {
		Id                    int         `json:"id"`
		DeptId                int         `json:"deptId"`
		Phone                 string      `json:"phone"`
		Avatar                interface{} `json:"avatar"`
		TenantId              int         `json:"tenantId"`
		Username              string      `json:"username"`
		NickName              string      `json:"nickName"`
		Password              interface{} `json:"password"`
		Enabled               bool        `json:"enabled"`
		AccountNonExpired     bool        `json:"accountNonExpired"`
		CredentialsNonExpired bool        `json:"credentialsNonExpired"`
		AccountNonLocked      bool        `json:"accountNonLocked"`
		Authorities           []struct {
			Authority string `json:"authority"`
		} `json:"authorities"`
	} `json:"user_info"`
	UserName         string   `json:"user_name"`
	Scope            []string `json:"scope"`
	Active           bool     `json:"active"`
	Exp              string   `json:"exp"`
	Authorities      []string `json:"authorities"`
	ClientId         string   `json:"client_id"`
	Error            string   `json:"error"`
	ErrorDescription string   `json:"error_description"`
}

func Server(basicUsername, basicPassword, checkUrl string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				auths := strings.SplitN(header.RequestHeader().Get(AuthorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], BearerWord) {
					//return nil, ErrMissingJwtToken
					return handler(ctx, req)
				}
				token := auths[1]
				// handle oauth logic
				response, err := verifyToken(token, basicUsername, basicPassword, checkUrl)
				if err != nil {
					//return nil, ErrMissingJwtToken
					return handler(ctx, req)
				}
				ctx = NewContext(ctx, response)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

func verifyToken(token, apiKey, password, checkUrl string) (*Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?token=%s", checkUrl, token), nil)
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(apiKey, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrMissingJwtToken
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrMissingJwtToken
	}
	checkResp := Response{}
	err = json.Unmarshal(body, &checkResp)
	if err != nil {
		return nil, ErrMissingJwtToken
	}
	if len(checkResp.Error) != 0 || checkResp.Error == "invalid_token" {
		return nil, ErrMissingJwtToken
	}
	return &checkResp, nil
}

func NewContext(ctx context.Context, info *Response) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}
