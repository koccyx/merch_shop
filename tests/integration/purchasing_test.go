package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	_ "github.com/golang-jwt/jwt/v5"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PurchasingTests struct {
	name string
	username string
	password string
	item string
	res int
	authRes int
}

func TestPurchase(t *testing.T) {
	tests := []PurchasingTests{
		{
			name: "Success test",
			username: RandomWord(7),
			password: RandomWord(7),
			item: "socks",
			res: http.StatusOK,
			authRes: http.StatusOK,
		},
		{
			name: "Failure item name test",
			username:     RandomWord(7),
			password:     RandomWord(7),
			item: "pencil",
			res: http.StatusBadRequest,
			authRes: http.StatusOK,
		},
		{
			name: "Failure auth test",
			username:     RandomWord(5),
			password:     RandomWord(3),
			item: "pen",
			res: http.StatusUnauthorized,
			authRes: http.StatusBadRequest,
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
			assert.Equal(t, tt.authRes, res.StatusCode)


			if (tt.authRes == http.StatusOK) {
				var authResponse = models.AuthResponse{}
				err = json.NewDecoder(res.Body).Decode(&authResponse)
				require.NoError(t, err)

				purReq, err := http.NewRequest("GET", apiURL+"/api/buy/"+tt.item, nil)
				require.NoError(t, err)
				
				purReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.Token))
				purReq.Header.Set("Content-Type", "application/json")

				purRes, err := client.Do(purReq)
				
				require.NoError(t, err)
				defer purRes.Body.Close()

				assert.Equal(t, tt.res, purRes.StatusCode)
			}
        })
	}
}