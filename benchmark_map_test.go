package fgo_test

import (
	"fmt"
	"testing"

	f "github.com/slavsan/fgo"
	"github.com/slavsan/fgo/internal/assert"
)

func imperativeMap(input []string) []string {
	output := make([]string, 0, len(input))
	for _, v := range input {
		output = append(output, v+"2")
	}
	return output
}

func funMap(t *testing.T, input []string) []string {
	output, _, err :=
		f.Pipe[string, string, any](
			input,
			f.Map(func(s string) string { return s + "2" }),
		)
	if t != nil {
		assert.Nil(t, err)
	}

	return output
}

func TestMapImplementations(t *testing.T) {
	input := []string{"foo", "bar", "baz"}
	expected := []string{"foo2", "bar2", "baz2"}
	assert.Equal(t, expected, imperativeMap(input))
	assert.Equal(t, expected, funMap(t, input))
}

func BenchmarkImperativeMap(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = imperativeMap(input)
			}
		})
	}
}

func BenchmarkFunctionalMap(b *testing.B) {
	testCases := []int{10, 100, 1000, 10000}
	for _, tc := range testCases {
		tc := tc
		b.Run(fmt.Sprintf("with %d items", tc), func(b *testing.B) {
			input := generateStrings(tc)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = funMap(nil, input)
			}
		})
	}
}
