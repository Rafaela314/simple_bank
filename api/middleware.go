package api

import (
	"errors"
	"fmt"
	"net/http"
	"simple_bank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

var (
	errMissingAuthHeader = errors.New("authorization header is required")
	errInvalidAuthFormat = errors.New("invalid authorization header format")
	errEmptyToken        = errors.New("authorization token is empty")
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		payload, err := extractBearerToken(ctx, tokenMaker)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// extractBearerToken parses the Authorization header and verifies the Bearer token.
func extractBearerToken(ctx *gin.Context, tokenMaker token.Maker) (*token.Payload, error) {
	authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		return nil, errMissingAuthHeader
	}
	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return nil, errInvalidAuthFormat
	}
	if strings.ToLower(fields[0]) != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", fields[0])
	}
	accessToken := fields[1]
	if accessToken == "" {
		return nil, errEmptyToken
	}
	return tokenMaker.VerifyToken(accessToken)
}

// GetAuthorizationPayload returns the token payload set by authMiddleware.
// Use this in handlers to get the authenticated user's payload. Returns (nil, false) if not set.
func GetAuthorizationPayload(ctx *gin.Context) (*token.Payload, bool) {
	v, exists := ctx.Get(authorizationPayloadKey)
	if !exists || v == nil {
		return nil, false
	}
	payload, ok := v.(*token.Payload)
	return payload, ok
}
