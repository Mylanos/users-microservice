package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"users-backend/pkg/api"
	"users-backend/pkg/config"
	"users-backend/pkg/services"
	"users-backend/pkg/storage"
)

const (
	testDatabaseURL = "postgresql://test_user:test_password@localhost:5433/test_db?sslmode=disable"
)

type TestSuite struct {
	storage *storage.PostgresStorage
	service services.UserService
	server  *api.APIServer
	httpSrv *httptest.Server
	Client  *http.Client
}

func SetupTestSuite(t *testing.T) *TestSuite {
	cfg := &config.Config{
		DatabaseURL: testDatabaseURL,
	}

	testStorage, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		t.Fatalf("FATAL: failed to create a storage: %s", err)
	}

	testService, err := services.NewUserService(testStorage)
	if err != nil {
		t.Fatalf("FATAL: failed to create test service: %v", err)
	}

	apiServer := api.NewAPIServer(":8081", testService)

	httpServer := httptest.NewServer(apiServer.Handler())
	client := &http.Client{Timeout: 10 * time.Second}

	return &TestSuite{
		storage: testStorage,
		service: testService,
		server:  apiServer,
		httpSrv: httpServer,
		Client:  client,
	}
}

// cleans up the test environment
func (ts *TestSuite) Teardown(t *testing.T) {
	if ts.httpSrv != nil {
		ts.httpSrv.Close()
	}

	if err := ts.storage.CleanupTable(); err != nil {
		t.Errorf("Failed to clean up table: %v", err)
	}

	if ts.storage != nil {
		if err := ts.storage.Close(); err != nil {
			t.Errorf("Failed to close database connection: %v", err)
		}
	}
}

func (ts *TestSuite) makeJSONRequest(t *testing.T, method, url string, payload interface{}) *http.Response {
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Errorf("Failed to encode payload: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ts.Client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	return resp
}

func (ts *TestSuite) makeGETRequest(t *testing.T, url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	resp, err := ts.Client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	log.Printf("📥 GET Response: %d %s", resp.StatusCode, resp.Status)
	return resp
}
