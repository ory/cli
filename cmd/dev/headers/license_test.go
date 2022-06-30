package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYmlComment(t *testing.T) {
	tests := map[string]string{
		"Hello":        "# Hello\n",          // single line text
		"Hello\nWorld": "# Hello\n# World\n", // multi-line text
	}
	for give, want := range tests {
		t.Run(fmt.Sprintf("%s -> %s", give, want), func(t *testing.T) {
			t.Parallel()
			have := ymlComment(give)
			assert.Equal(t, want, have)
		})
	}
}

func TestRemoveHeader(t *testing.T) {
	t.Parallel()
	give := `# Copyright © 1997 Ory Corp Inc.

name: test
hello: world`
	want := `name: test
hello: world`
	have := removeHeader(give, ymlComment, LICENSE_TOKEN)
	assert.Equal(t, want, have)
}
