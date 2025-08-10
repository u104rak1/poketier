package str_test

import (
	"poketier/pkg/str"
	"reflect"
	"testing"
)

func TestCommaSeparatedToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "通常のカンマ区切り",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "スペースを含む",
			input:    " a, b ,c ",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "空文字列",
			input:    "",
			expected: []string{},
		},
		{
			name:     "カンマのみ",
			input:    ",,,",
			expected: []string{},
		},
		{
			name:     "値と空要素混在",
			input:    "a,,b, ,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "前後にスペースとカンマ",
			input:    "  , a , ",
			expected: []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.CommaSeparatedToSlice(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("input=%q: want %v, got %v", tt.input, tt.expected, got)
			}
		})
	}
}
