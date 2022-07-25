package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function that indicates within a test that we create a test file
// and encapsulates incidental complexity.
func CreateTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0744)
	assert.NoError(t, err)
}

func Trim(text string) string {
	return strings.Trim(text, "\n")
}
