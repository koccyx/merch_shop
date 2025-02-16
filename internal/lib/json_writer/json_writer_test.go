package jsonwriter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success test",
			input:          map[string]string{"id": "12312"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"12312"}`,
		},
		{
			name:           "empty input case",
			input:          map[string]string{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := WriteJSON(tt.input, rr)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		status         int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "bad request case",
			inputError:     errors.New("something went wrong"),
			status:         http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := WriteJSONError(tt.inputError, rr, tt.status)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestReadReqJSON(t *testing.T) {
	tests := []struct {
		name          string
		inputBody     string
		expectedData  map[string]string
		expectedError bool
	}{
		{
			name:          "success test",
			inputBody:     `{"key": "value"}`,
			expectedData:  map[string]string{"key": "value"},
			expectedError: false,
		},
		{
			name:          "json with no key fail",
			inputBody:     `{"key":}`,
			expectedData:  nil,
			expectedError: true,
		},
		{
			name:          "many fields",
			inputBody:     `{"key": "value", "extra": "field"}`,
			expectedData:  map[string]string{"key": "value", "extra": "field"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			var data map[string]string

			err := ReadReqJSON(rr, req, &data)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		rr := httptest.NewRecorder()

		WriteSuccess(rr)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}
