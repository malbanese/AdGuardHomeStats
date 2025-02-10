package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	StubUsername = "User"
	StubPassword = "Password"
)

func createThrowingServer(errorCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorCode)
	}))
}

func createSuccessServer(response *AghStatsResponse) *httptest.Server {
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
		client := http.Client{}

		expected := AghStatsResponse{
			TimeUnits:           "hours",
			NumDnsQueries:       1000,
			NumBlockedFiltering: 2000,
			AvgProcessingTime:   0.123,
		}

		server := createSuccessServer(&expected)
		defer server.Close()

		endpoint := AghStatsEndpointInfo{
			Username: StubUsername,
			Password: StubPassword,
			Url:      server.URL,
		}

		var response AghStatsResponse
		err := FetchStats(&client, &response, &endpoint)

		if err != nil {
			t.Errorf("Expected success, got error instead %v", err)
		} else if response != expected {
			t.Error("Expected response to equal initial")
		}
	})
}

func TestFetchStatsErrors(t *testing.T) {
	testCases := []struct {
		name        string
		errorPrefix string
		endpoint    *AghStatsEndpointInfo
		server      *httptest.Server
	}{
		{
			name:        "Error making request object",
			errorPrefix: "error creating request",
			endpoint: &AghStatsEndpointInfo{
				Username: StubUsername,
				Password: StubPassword,
				Url:      "http://user:abcd{DEf1=ghi@example.com:5432/db",
			},
			server: createSuccessServer(nil),
		},
		{
			name:        "Error executing request",
			errorPrefix: "error making request",
			endpoint: &AghStatsEndpointInfo{
				Username: StubUsername,
				Password: StubPassword,
				Url:      "slash_at_end_is_invalid/",
			},
			server: createSuccessServer(nil),
		},
		{
			name:        "Downstream responds with 400",
			errorPrefix: "error parsing body json",
			endpoint: &AghStatsEndpointInfo{
				Username: StubUsername,
				Password: StubPassword,
				Url:      "",
			},
			server: createSuccessServerWithString("{ invalid json"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.endpoint.Url == "" {
				tc.endpoint.Url = tc.server.URL
			}

			defer tc.server.Close()
			var response AghStatsResponse
			client := http.Client{}
			err := FetchStats(&client, &response, tc.endpoint)
			if err == nil {
				t.Errorf("Expected throwing case for `%s`", tc.name)
			} else if !strings.HasPrefix(err.Error(), tc.errorPrefix) {
				t.Errorf("Expected prefix `%s` for `%v` on `%s`", tc.errorPrefix, err, tc.name)
			}
		})
	}
}
