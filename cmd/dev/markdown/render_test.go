package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	assert.EqualValues(t, `<strong>foo</strong>`, renderMarkdown([]byte(`**foo**`)))
}
