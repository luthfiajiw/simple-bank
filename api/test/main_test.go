package api_test

import (
	"os"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *api.Server {
	config := utils.Config{
		SymetricKey:         "1eb9dbbbbc047c03fd70604e0071f098",
		AccessTokenDuration: time.Minute,
	}

	server, err := api.NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
