package casdoor_auth

import (
	"context"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"strings"
)

type AuthKey struct{}

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

// NewWhiteListMatcher 白名单过滤
func NewWhiteListMatcher(ignoreUrl []string) selector.MatchFunc {
	whiteList := make(map[string]struct{})
	for _, url := range ignoreUrl {
		whiteList[url] = struct{}{}
	}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// Server 认证服务
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				auths := strings.SplitN(header.RequestHeader().Get(AuthorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], BearerWord) {

					return nil, ErrMissingJwtToken
				}
				token := auths[1]
				response, err := verifyToken(token)
				if err != nil {
					return handler(ctx, req)
				}
				ctx = NewContext(ctx, response)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client 声明 casdoor 的初始化
func Client(endpoint string, clientId string, clientSecret string, jwtPublicKey string, organizationName string, applicationName string) {
	auth.InitConfig(endpoint, clientId, clientSecret, jwtPublicKey, organizationName, applicationName)
}

func verifyToken(token string) (*auth.Claims, error) {
	claims, err := auth.ParseJwtToken(token)
	if err != nil {
		return nil, ErrMissingJwtToken
	}

	return claims, nil
}

func NewContext(ctx context.Context, info *auth.Claims) context.Context {
	return context.WithValue(ctx, AuthKey{}, info)
}

func GetUserInfo(ctx context.Context) (*auth.Claims, bool) {
	info, ok := ctx.Value(AuthKey{}).(*auth.Claims)
	return info, ok
}
