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

type SendingCoinsTests struct {
	name string
	username1 string
	password1 string
	username2 string
	password2 string
	amount int
	res int
	authRes int
}

func TestSending(t *testing.T) {
	tests := []SendingCoinsTests{
		{
			name: "Success test",
			username1: RandomWord(10),
			password1: RandomWord(10),
			username2: RandomWord(10),
			password2: RandomWord(10),
			amount: 100,
			res: http.StatusOK,
			authRes: http.StatusOK,
		},
		// {
		// 	name: "Failure item name test",
		// 	username:     RandomWord(7),
		// 	password:     RandomWord(7),
		// 	item: "pencil",
		// 	res: http.StatusBadRequest,
		// 	authRes: http.StatusOK,
		// },
		// {
		// 	name: "Failure auth test",
		// 	username:     RandomWord(5),
		// 	password:     RandomWord(3),
		// 	item: "pen",
		// 	res: http.StatusUnauthorized,
		// 	authRes: http.StatusBadRequest,
		// },
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}

			reqBody1 := &models.AuthRequest{
				Username: tt.username1,
				Password: tt.password1,
			}

			r1, err := json.Marshal(reqBody1)
			require.NoError(t, err)

			req1, err := http.NewRequest("POST", apiURL+"/api/auth", bytes.NewBuffer(r1))
			require.NoError(t, err)
			
			req1.Header.Set("Content-Type", "application/json")
			res1, err := client.Do(req1)
			require.NoError(t, err)
			defer res1.Body.Close()
			assert.Equal(t, tt.authRes, res1.StatusCode)


			reqBody2 := &models.AuthRequest{
				Username: tt.username1,
				Password: tt.password1,
			}

			r2, err := json.Marshal(reqBody2)
			require.NoError(t, err)

			req2, err := http.NewRequest("POST", apiURL+"/api/auth", bytes.NewBuffer(r2))
			require.NoError(t, err)
			
			req2.Header.Set("Content-Type", "application/json")
			res2, err := client.Do(req2)
			require.NoError(t, err)

			defer res2.Body.Close()
			assert.Equal(t, tt.authRes, res2.StatusCode)

			if (tt.authRes == http.StatusOK) {
				authResponse1 := models.AuthResponse{}
				err = json.NewDecoder(res1.Body).Decode(&authResponse1)
				require.NoError(t, err)

				sendReq := models.SendCoinRequest{
					ToUser: tt.username2,
					Amount: tt.amount,
				}

				sendR, err := json.Marshal(sendReq)
				require.NoError(t, err)

				purReq, err := http.NewRequest("post", apiURL+"/api/sendCoin", bytes.NewBuffer(sendR))
				require.NoError(t, err)
				
				purReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResponse1.Token))
				purReq.Header.Set("Content-Type", "application/json")

				purRes, err := client.Do(purReq)
				
				require.NoError(t, err)
				defer purRes.Body.Close()

				assert.Equal(t, tt.res, purRes.StatusCode)
			}
        })
	}
}