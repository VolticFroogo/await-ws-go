package awaitws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_IsRequest(t *testing.T) {
	client := NewClient(nil)

	tests := []struct {
		name     string
		arg      map[string]interface{}
		expected bool
	}{
		{
			name: "True",
			arg: map[string]interface{}{
				"request": true,
			},
			expected: true,
		},
		{
			name: "False",
			arg: map[string]interface{}{
				"request": false,
			},
			expected: false,
		},
		{
			name:     "Empty",
			arg:      map[string]interface{}{},
			expected: false,
		},
		{
			name: "Invalid type",
			arg: map[string]interface{}{
				"request": "badType",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, client.IsRequest(tt.arg))
		})
	}
}
