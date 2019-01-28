package mergo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type SomeStruct struct {
	Integer int `json:"int_tag,omitempty"`
	Boolean bool
}

func TestMapWithFieldTag(t *testing.T) {
	tests := []struct {
		name     string
		src      map[string]interface{}
		expected SomeStruct
	}{
		{
			name: "empty map, nil dst",
		},
		{
			name: "map integer with json field tag",
			src: map[string]interface{}{
				"int_tag": 5,
				"Boolean": true,
			},
			expected: SomeStruct{
				Integer: 5,
				Boolean: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := SomeStruct{}
			assert.NoError(t, Map(&dst, &tt.src, WithFieldTag("json")))
			assert.Equal(t, tt.expected, dst)
		})
	}
}
