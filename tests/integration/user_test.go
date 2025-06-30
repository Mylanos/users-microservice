package integration

import (
	"io"
	"log"
	"net/http"
	"testing"
	"time"
	"users-microservice/pkg/api"

	"github.com/google/uuid"
)

func TestUserEndpoints(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Teardown(t)

	testCasesPOST := []struct {
		name     string
		reqData  api.UserAPI
		wantCode int
	}{
		{
			name:     "valid user",
			reqData:  api.UserAPI{ID: uuid.New(), Name: "Milan", Email: "haha@google.com", DateOfBirth: time.Now().AddDate(-25, 0, 0)},
			wantCode: 201,
		},
		{
			name:     "invalid email",
			reqData:  api.UserAPI{ID: uuid.New(), Name: "Stefan", Email: "invalid", DateOfBirth: time.Now().AddDate(-25, 0, 0)},
			wantCode: 400,
		},
		{
			name:     "duplicate email",
			reqData:  api.UserAPI{ID: uuid.New(), Name: "Stefan", Email: "haha@google.com", DateOfBirth: time.Now().AddDate(-32, 0, 0)},
			wantCode: 409,
		},
		{
			name:     "future birthday",
			reqData:  api.UserAPI{ID: uuid.New(), Name: "Milada", Email: "new@google.com", DateOfBirth: time.Now().AddDate(32, 0, 0)},
			wantCode: 400,
		},
		{
			name:     "empty name",
			reqData:  api.UserAPI{ID: uuid.New(), Name: "", Email: "else@google.com", DateOfBirth: time.Now().AddDate(-76, 0, 0)},
			wantCode: 400,
		},
	}

	for _, tc := range testCasesPOST {
		t.Run(tc.name, func(t *testing.T) {
			resp := suite.makeJSONRequest(t, "POST", suite.httpSrv.URL+"/save", tc.reqData)
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantCode {
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Response for '%s': Status=%d, Body=%s", tc.name, resp.StatusCode, string(body))
				t.Errorf("Test '%s': Expected status %d, got %d", tc.name, tc.wantCode, resp.StatusCode)
			}
		})
	}

	testCasesGET := []struct {
		name     string
		reqID    string
		wantCode int
	}{
		{
			name:     "valid id",
			reqID:    uuid.New().String(),
			wantCode: 200,
		},
		{
			name:     "bad format id",
			reqID:    "bad",
			wantCode: 400,
		},
	}

	for _, tc := range testCasesGET {
		t.Run(tc.name, func(t *testing.T) {
			parsedUUID, err := uuid.Parse(tc.reqID)
			if err == nil {
				createReq := api.UserAPI{
					ID:          parsedUUID,
					Name:        "Milan",
					Email:       "milan@test.com",
					DateOfBirth: time.Now().AddDate(-25, 0, 0),
				}
				resp1 := suite.makeJSONRequest(t, "POST", suite.httpSrv.URL+"/save", createReq)
				defer resp1.Body.Close()

				if resp1.StatusCode != http.StatusCreated {
					body, _ := io.ReadAll(resp1.Body)
					t.Fatalf("Failed to create user: Status=%d, Body=%s", resp1.StatusCode, string(body))
				}
			}

			resp2 := suite.makeGETRequest(t, suite.httpSrv.URL+"/"+tc.reqID)
			defer resp2.Body.Close()

			if resp2.StatusCode != tc.wantCode {
				body, _ := io.ReadAll(resp2.Body)
				log.Printf("Response for '%s': Status=%d, Body=%s", tc.name, resp2.StatusCode, string(body))
				t.Errorf("Test '%s': Expected status %d, got %d", tc.name, tc.wantCode, resp2.StatusCode)
			}
		})
	}

	log.Print("attempting unknown methods")
	t.Run("unknown_method_on_get_endpoint", func(t *testing.T) {
		resp := suite.makeJSONRequest(t, "PUT", suite.httpSrv.URL+"/"+uuid.New().String(), "")
		defer resp.Body.Close()
		if resp.StatusCode != 405 {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Response for PUT method on GET: Status=%d, Body=%s", resp.StatusCode, string(body))
			t.Errorf("Test: Expected status %d, got %d", 405, resp.StatusCode)
		}
	})
	t.Run("unknown_method_on_post_endpoint", func(t *testing.T) {
		resp2 := suite.makeJSONRequest(t, "PUT", suite.httpSrv.URL+"/save", "")
		defer resp2.Body.Close()

		if resp2.StatusCode != 405 {
			body, _ := io.ReadAll(resp2.Body)
			log.Printf("Response for PUT method on GET: Status=%d, Body=%s", resp2.StatusCode, string(body))
			t.Errorf("Test: Expected status %d, got %d", 405, resp2.StatusCode)
		}
	})
}
