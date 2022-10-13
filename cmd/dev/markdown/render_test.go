// Copyright © 2022 Ory Corp

package markdown

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	assert.EqualValues(t, `<strong>foo</strong><br>`,
		strings.TrimSpace(string(renderMarkdown([]byte(`**foo**`)))))
}
