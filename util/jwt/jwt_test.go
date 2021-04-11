package jwt

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJWT(t *testing.T) {
	m := map[string]interface{}{
		"id":   "1",
		"name": "fenghai",
	}
	raw, err := Token(m)
	require.NoError(t, err)
	fmt.Println("token:", raw)
	om, err := ParseToken(raw)
	require.NoError(t, err)
	for i, v := range m {
		if ov, ok := om[i]; !ok || v != ov {
			panic("token parse error")
		}
	}
}
