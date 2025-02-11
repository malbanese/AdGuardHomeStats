package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createThrowingServer(errorCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorCode)
	}))
}

func createSuccessServer(response *StatsResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
		}
	}))
}

func createSuccessServerWithString(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, response)
	}))
}

func TestFetchStatsSuccess(t *testing.T) {
	t.Run("Fetch successful response", func(t *testing.T) {

		expected := StatsResponse{
			TimeUnits:           "hours",
			NumDnsQueries:       1000,
			NumBlockedFiltering: 2000,
			AvgProcessingTime:   0.123,
		}

		server := createSuccessServer(&expected)
		defer server.Close()

		client := StatsClient{
			Client: &http.Client{},
			Url:    server.URL,
			Auth: &BasicAuthType{
				Username: "user",
				Password: "password",
			},
		}

		var response StatsResponse
		err := client.FetchStats(&response)

		if err != nil {
			t.Errorf("Expected success, got error instead %v", err)
		} else if response != expected {
			t.Error("Expected response to equal initial")
		}
	})
}

func TestFetchStatsErrors(t *testing.T) {
	username := "username"
	password := "password"

	testCases := []struct {
		name        string
		errorPrefix string
		username    string
		password    string
		url         string
		server      *httptest.Server
	}{
		{
			name:        "Error making request object",
			errorPrefix: "error creating request",
			username:    username,
			password:    password,
			url:         "http://user:abcd{DEf1=ghi@example.com:5432/db",
			server:      createSuccessServer(nil),
		},
		{
			name:        "Error executing request",
			errorPrefix: "error making request",
			username:    username,
			password:    password,
			url:         "slash_at_end_is_invalid/",
			server:      createSuccessServer(nil),
		},
		{
			name:        "Downstream responds with 400",
			errorPrefix: "HTTP Error",
			username:    username,
			password:    password,
			url:         "",
			server:      createThrowingServer(400),
		},
		{
			name:        "Downstream gives invalid json",
			errorPrefix: "error parsing body json",
			username:    username,
			password:    password,
			url:         "",
			server:      createSuccessServerWithString("{ invalid json"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.url == "" {
				tc.url = tc.server.URL
			}

			defer tc.server.Close()
			client := StatsClient{
				Client: &http.Client{},
				Url:    tc.url,
				Auth: &BasicAuthType{
					Username: "user",
					Password: "password",
				},
			}

			var response StatsResponse
			err := client.FetchStats(&response)

			if err == nil {
				t.Errorf("Expected throwing case for `%s`", tc.name)
			} else if !strings.HasPrefix(err.Error(), tc.errorPrefix) {
				t.Errorf("Expected prefix `%s` for `%v` on `%s`", tc.errorPrefix, err, tc.name)
			}
		})
	}
}
