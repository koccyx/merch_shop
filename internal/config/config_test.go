package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Success(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", "test_config.yml")
	require.NoError(t, err)

	configContent := `
env: test
db:
  addres: localhost
  port: "5432"
  database: testdb
  user: testuser
  password: testpass
  schema: public
http_server:
  addres: localhost
  port: "8080"
  timeout: 5s
  idle_timeout: 60s
auth:
  secret: supersecret`

	file, err := os.Create("test_config.yml")
	require.NoError(t, err)
	defer os.Remove("test_config.yml")

	_, err = file.WriteString(configContent)
	require.NoError(t, err)

	file.Close()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "localhost", cfg.Storage.Addres)
	assert.Equal(t, "5432", cfg.Storage.Port)
	assert.Equal(t, "testdb", cfg.Storage.Database)
	assert.Equal(t, "testuser", cfg.Storage.User)
	assert.Equal(t, "testpass", cfg.Storage.Password)
	assert.Equal(t, "public", cfg.Storage.Schema)
	assert.Equal(t, "localhost", cfg.Server.Addres)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, 5*time.Second, cfg.Server.TimeOut)
	assert.Equal(t, 60*time.Second, cfg.Server.IdleTimeout)
	assert.Equal(t, "supersecret", cfg.Auth.Secret)
}

func TestLoadConfig_MissingConfigPath(t *testing.T) {
	err := os.Unsetenv("CONFIG_PATH")
	require.NoError(t, err)

	_, err = Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "CONFIG_PATH env variable not set")
}

func TestLoadConfig_FileDoesNotExist(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", "non_existent_config.yml")
	require.NoError(t, err)

	_, err = Load()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "config file doesnt exist")
}
