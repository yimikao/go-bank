package api

import (
	"fmt"
	"gobank/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorizationHeader(
	t *testing.T,
	r *http.Request,
	tm token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tm.CreateToken(username, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authorizationType, token)
	r.Header.Set(authorizationHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, r *http.Request, tm token.Maker)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, r *http.Request, tm token.Maker) {
				addAuthorizationHeader(t, r, tm, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "No authorization",
			setupAuth: func(t *testing.T, r *http.Request, tm token.Maker) {
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "Unsupported authorization",
			setupAuth: func(t *testing.T, r *http.Request, tm token.Maker) {
				addAuthorizationHeader(t, r, tm, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "Invalid authorization format",
			setupAuth: func(t *testing.T, r *http.Request, tm token.Maker) {
				addAuthorizationHeader(t, r, tm, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "Expired token",
			setupAuth: func(t *testing.T, r *http.Request, tm token.Maker) {
				addAuthorizationHeader(t, r, tm, "unsupported", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			path := "/auth"
			server.router.GET(
				path,
				authMiddleware(server.tokenMaker),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})
				},
			)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}
