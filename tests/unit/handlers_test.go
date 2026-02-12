package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortCode(t *testing.T) {
    // Test short code generation
    code := generateTestCode(6)
    assert.Equal(t, 6, len(code))
}

func generateTestCode(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    code := make([]byte, length)
    for i := range code {
        code[i] = charset[i%len(charset)]
    }
    return string(code)
}