package remote_test

import (
	"github.com/ory/cli/cmd/cloud/remote"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	InvalidAPIEndpoint = ""
	ValidAPIEndpoint = "https://oryapis:8080"
	InvalidConsoleURL = ""
	ValidConsoleURL  = "https://console.ory:8080"
)

var FlagTest = []struct{
	description string
	in []string
	out []string
}{
	{"Missing protocols",[]string{"oryapis:8080", "console.ory:8080"}, []string{"", "Could not parse requested url: console.ory:8080"}},
	{"Invalid urls",[]string{"hi/there?", "hi/there?"}, []string{"","invalid URI for request"}},
}

func TestClient(t *testing.T) {
	t.Run("function=GetProjectSlug", func(t *testing.T) {
		t.Run("should fail with invalid urls", func(t *testing.T) {
			for _, tt := range FlagTest {
				t.Run(tt.description, func(t *testing.T) {
					slug, err := remote.GetProjectSlug(tt.in[1])
					assert.Contains(t, err.Error(), tt.out[1])
					assert.EqualValues(t, tt.out[0], slug)
				})
			}
		})
	})

}
