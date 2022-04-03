package object

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString_Len(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"12345", 5},
		{"", 0},
		{"ffffffffff", 10},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			s := &String{Value: subtest.input}
			require.Equal(t, &Integer{Value: subtest.expected}, s.Len())
		})
	}
}
