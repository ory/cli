package newsletter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	assert.EqualValues(t, `<img width="600" border="0" style="display: block; max-width: 100%; min-width: 100px; width: 100%;" src=".foo/bar" alt="Some image" /><br>

<pre style="word-break: break-all; white-space: pre-wrap"><code>foo
</code></pre>`,
		strings.TrimSpace(string(renderMarkdown([]byte(`
![Some image](.foo/bar)

`+"```\nfoo\n```"+`

`)))))
}
