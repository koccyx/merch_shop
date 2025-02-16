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

type InfoTests struct {
	name string
	username1 string
	password1 string
	username2 string
	password2 string
	amount int
	res int
}

func TestInfo(t *testing.T) {
	tests := []SendingCoinsTests{
		{
			name: "Success test",
			username1: RandomWord(10),
			password1: RandomWord(10),
			username2: RandomWord(10),
			password2: RandomWord(10),
			amount: 100,
			res: http.StatusOK,
		},
		{
			name: "Failure not enough balance test",
			username1: RandomWord(10),
			password1: RandomWord(10),
			username2: RandomWord(10),
			password2: RandomWord(10),
			amount: 1001,
			res: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}

			reqBody := &models.AuthRequest{
				Username: tt.username1,
				Password: tt.password1,
			}

			r, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", apiURL+"/api/auth", bytes.NewBuffer(r))
			require.NoError(t, err)
			
			req.Header.Set("Content-Type", "application/json")
			res1, err := client.Do(req)
			require.NoError(t, err)
			defer res1.Body.Close()


			reqBody = &models.AuthRequest{
				Username: tt.username2,
				Password: tt.password2,
			}

			r, err = json.Marshal(reqBody)
			require.NoError(t, err)

			req, err = http.NewRequest("POST", apiURL+"/api/auth", bytes.NewBuffer(r))
			require.NoError(t, err)
			
			req.Header.Set("Content-Type", "application/json")
			res2, err := client.Do(req)
			require.NoError(t, err)
			defer res2.Body.Close()

			authResponse := models.AuthResponse{}
			err = json.NewDecoder(res1.Body).Decode(&authResponse)
			require.NoError(t, err)
		
			sendReq := models.SendCoinRequest{
				ToUser: tt.username2,
				Amount: tt.amount,
			}

			sendR, err := json.Marshal(sendReq)
			require.NoError(t, err)

			purReq, err := http.NewRequest("POST", apiURL+"/api/sendCoin", bytes.NewBuffer(sendR))
			require.NoError(t, err)

			purReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.Token))
			purReq.Header.Set("Content-Type", "application/json")

			purRes, err := client.Do(purReq)
			require.NoError(t, err)
			defer purRes.Body.Close()

			assert.Equal(t, tt.res, purRes.StatusCode)

			infReq, err := http.NewRequest("GET", apiURL+"/api/info", nil)
			require.NoError(t, err)

			infReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.Token))
			infReq.Header.Set("Content-Type", "application/json")

			infRes, err := client.Do(infReq)
			require.NoError(t, err)
			defer infRes.Body.Close()

			infoResp := models.InfoResponse{}
			err = json.NewDecoder(infRes.Body).Decode(&infoResp)
			require.NoError(t, err)
			

			if tt.res == http.StatusOK {
				assert.Equal(t, 1000-tt.amount, infoResp.Coins)
			}
        })
	}
}