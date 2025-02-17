package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simplebank/api"
	"simplebank/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	authType string,
	username string,
	durarion time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, durarion)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, token)
	req.Header.Set(api.AuthorizationHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuth(t, req, tokenMaker, api.AuthorizationTypeBearer, "luthfi", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.Router.GET(
				authPath,
				api.AuthMiddleware(server.TokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.TokenMaker)
			server.Router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
