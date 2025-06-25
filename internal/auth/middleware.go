package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

var (
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidAuthorizationType   = errors.New("invalid authorization type")
	ErrInvalidAuthorizationFormat = errors.New("invalid authorization format")
)

// BearerMiddleware returns a Gin middleware for Bearer token authentication.
func BearerMiddleware(tokenMaker TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := strings.TrimSpace(ctx.GetHeader(authorizationHeaderKey))
		if header == "" {
			httpUnauthorized(ctx, ErrInvalidAuthorizationHeader)
			return
		}

		fields := strings.Fields(header)
		if len(fields) != 2 {
			httpUnauthorized(ctx, ErrInvalidAuthorizationFormat)
			return
		}

		if !strings.EqualFold(fields[0], authorizationTypeBearer) {
			httpUnauthorized(ctx, ErrInvalidAuthorizationType)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			httpUnauthorized(ctx, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func httpUnauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error":  err.Error(),
		"status": http.StatusUnauthorized,
	})
}
