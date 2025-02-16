package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/koccyx/avito_assignment/internal/lib/jwt"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthTests struct {
	name     string
	username string
	password string
	res      int
}

func TestAuth(t *testing.T) {
	tests := []AuthTests{
		{
			name:     "Success test",
			username: RandomWord(7),
			password: RandomWord(7),
			res:      http.StatusOK,
		},
		{
			name:     "Failure short password test",
			username: RandomWord(7),
			password: RandomWord(4),
			res:      http.StatusBadRequest,
		},
		{
			name:     "Failure short username test",
			username: RandomWord(3),
			password: RandomWord(7),
			res:      http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}

			reqBody := &models.AuthRequest{
				Username: tt.username,
				Password: tt.password,
			}

			r, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", apiURL+"/api/auth", bytes.NewBuffer(r))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.res, res.StatusCode)

			if tt.res == http.StatusOK {
				var authResponse = models.AuthResponse{}
				err = json.NewDecoder(res.Body).Decode(&authResponse)
				require.NoError(t, err)

				_, err = jwt.ParseToken(authResponse.Token, cfg.Auth.Secret)
				require.NoError(t, err)
			}
		})
	}
}
