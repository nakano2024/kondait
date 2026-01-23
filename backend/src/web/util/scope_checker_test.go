package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAnyScope(t *testing.T) {
	testTable := []struct {
		name      string
		requireds []string
		granteds  []string
		expected  bool
	}{
		{
			name:      "必要スコープに一致がある場合、trueを返すこと",
			requireds: []string{"scope:read", "scope:write"},
			granteds:  []string{"scope:write"},
			expected:  true,
		},
		{
			name:      "必要スコープに一致がない場合、falseを返すこと",
			requireds: []string{"scope:read"},
			granteds:  []string{"scope:admin"},
			expected:  false,
		},
		{
			name:      "requiredsが空の場合、falseを返すこと",
			requireds: []string{},
			granteds:  []string{"scope:read"},
			expected:  false,
		},
		{
			name:      "grantedsが空の場合、falseを返すこと",
			requireds: []string{"scope:read"},
			granteds:  []string{},
			expected:  false,
		},
		{
			name:      "requiredsとgrantedsが空の場合、falseを返すこと",
			requireds: []string{},
			granteds:  []string{},
			expected:  false,
		},
		{
			name:      "requiredsに同一スコープが含まれる場合、trueを返すこと",
			requireds: []string{"scope:read", "scope:read"},
			granteds:  []string{"scope:read"},
			expected:  true,
		},
		{
			name:      "空文字スコープが一致する場合、trueを返すこと",
			requireds: []string{""},
			granteds:  []string{""},
			expected:  true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			got := HasAnyScope(tt.requireds, tt.granteds)
			assert.Equal(t, tt.expected, got)
		})
	}
}
