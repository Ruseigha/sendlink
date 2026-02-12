package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
    resp, err := http.Get("http://localhost:8080/health")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAndRedirectURL(t *testing.T) {
    // Create short URL
    payload := map[string]string{"url": "https://google.com"}
    jsonData, _ := json.Marshal(payload)
    
    resp, err := http.Post(
        "http://localhost:8080/api/shorten",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}