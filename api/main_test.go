package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/thanhphuocnguyen/go-simple-bank/db/sqlc"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

func createNewServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
