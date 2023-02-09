// 生成csrf_token
package encrypt

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
)

// 产生token
func GenCSRFToken() string {
	var buf bytes.Buffer
	buf.Grow(32)
	for i := 0; i < 8; i++ {
		a := 65536.0 * (1.0 + rand.Float64())
		b := int(math.Floor(a))
		c := fmt.Sprintf("%x", b)
		buf.WriteString(c[1:])
	}
	return buf.String()
}
