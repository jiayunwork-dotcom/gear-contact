package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHealthAndMesh(t *testing.T) {
	root := testdataRoot(t)
	h := New(filepath.Join(root, "web"), filepath.Join(root, "example", "spur-20.json"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health %d", rr.Code)
	}
	body, err := os.ReadFile(filepath.Join(root, "example", "spur-20.json"))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mesh", bytes.NewReader(body))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("mesh %d %s", rr.Code, rr.Body.String())
	}
	var got MeshResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ContactRatio <= 1 {
		t.Fatalf("epsilon %v", got.ContactRatio)
	}
	if got.CenterDistance != 30 {
		t.Fatalf("a %v", got.CenterDistance)
	}
}

func testdataRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "example", "spur-20.json")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("cannot find example")
	return ""
}
