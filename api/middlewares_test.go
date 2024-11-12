package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/thanhphuocnguyen/go-simple-bank/auth"
)

func addAuthHeader(
	t *testing.T,
	request *http.Request,
	tokenGenerator auth.TokenGenerator,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenGenerator.GenerateToken(username, duration)
	require.NoError(t, err)

	tokeBearer := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorization, tokeBearer)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name            string
		setupAuthHeader func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, "thanh", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				// Do nothing
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, "Token", "thanh", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidFormatAuthorization",
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, "", "thanh", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuthHeader: func(t *testing.T, request *http.Request, tokenGenerator auth.TokenGenerator) {
				addAuthHeader(t, request, tokenGenerator, authorizationType, "thanh", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	server := createNewServer(t, nil)
	authPath := "/auth"
	server.router.GET(authPath, authMiddleware(server.tokenGenerator), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuthHeader(t, req, server.tokenGenerator)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
